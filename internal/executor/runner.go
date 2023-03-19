package executor

import (
	"encoding/json"
	"errors"
	"os"
	"serverless/internal/engine"
	"serverless/internal/engine/functions/http"
	"serverless/internal/utils/file"
)

func ParseManifest(folder string) (*engine.Manifest, error) {
	var manifestJson, err = os.ReadFile(folder + "/manifest.json")
	if err != nil {
		return &engine.Manifest{}, err
	}
	var manifest engine.Manifest
	err = json.Unmarshal(manifestJson, &manifest)
	if err != nil {
		return &engine.Manifest{}, err
	}
	return &manifest, nil
}

func Execute(folder string, manifest *engine.Manifest) (chan http.Response, error) {
	fileName := folder + "/" + manifest.Main
	if !file.IsTextFile(fileName) {
		return nil, errors.New("invalid main file")
	}
	mainContent, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	e := engine.New(manifest)
	return e.RunScript(string(mainContent), fileName)
}
