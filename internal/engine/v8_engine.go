package engine

import (
	"fmt"
	"rogchap.com/v8go"
)

type Engine struct {
	context  *v8go.Context
	version  string
	name     string
	packages []string
}

type Manifest struct {
	Main     string   `json:"main"`
	Name     string   `json:"name"`
	Version  string   `json:"version"`
	BuildDir string   `json:"build_dir"`
	Packages []string `json:"packages"`
}

func registerFunctionCallbacks(isolate *v8go.Isolate, folder string) *v8go.ObjectTemplate {
	importTemplate := v8go.NewFunctionTemplate(isolate, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		fmt.Printf("%v", info.Args())
		return nil
	})
	global := v8go.NewObjectTemplate(isolate)
	global.Set("print", importTemplate)
	return global
}

func New(manifest *Manifest, folder string) *Engine {
	isolate := v8go.NewIsolate()
	callbacks := registerFunctionCallbacks(isolate, folder)
	context := v8go.NewContext(isolate, callbacks)

	return &Engine{
		context:  context,
		version:  manifest.Version,
		name:     manifest.Name,
		packages: manifest.Packages,
	}
}

func (e *Engine) RunScript(source string, file string) (*v8go.Value, error) {
	return e.context.RunScript(source, file)
}
