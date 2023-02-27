package engine

import (
	"rogchap.com/v8go"
	"serveless/internal/engine/functions/http"
)

type Engine struct {
	version  string
	name     string
	packages []string
	isolate  *v8go.Isolate
}

type Manifest struct {
	Main     string   `json:"main"`
	Name     string   `json:"name"`
	Version  string   `json:"version"`
	BuildDir string   `json:"build_dir"`
	Packages []string `json:"packages"`
}

func registerFunctionCallbacks(isolate *v8go.Isolate, channel chan http.Response) *v8go.ObjectTemplate {
	global := v8go.NewObjectTemplate(isolate)
	global.Set("response", http.ResponseFunction(isolate, channel))
	return global
}

func New(manifest *Manifest) *Engine {
	isolate := v8go.NewIsolate()
	return &Engine{
		isolate:  isolate,
		version:  manifest.Version,
		name:     manifest.Name,
		packages: manifest.Packages,
	}
}

func (e *Engine) CreateContext() (*v8go.Context, chan http.Response) {
	channel := make(chan http.Response, 1)
	callbacks := registerFunctionCallbacks(e.isolate, channel)
	context := v8go.NewContext(e.isolate, callbacks)
	return context, channel
}

func (e *Engine) RunScript(source string, file string) (chan http.Response, error) {
	context, responseChannel := e.CreateContext()
	go func() {
		_, err := context.RunScript(source, file)
		if err != nil {
			return
		}
	}()
	return responseChannel, nil
}
