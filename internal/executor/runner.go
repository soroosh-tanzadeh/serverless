package executor

import (
	"encoding/json"
	"errors"
	"os"
	"rogchap.com/v8go"
	"serveless/internal/engine"
	"serveless/internal/utils/file"
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

func Execute(folder string, manifest *engine.Manifest) (*v8go.Value, error) {
	fileName := folder + "/" + manifest.Main
	if !file.IsTextFile(fileName) {
		return &v8go.Value{}, errors.New("invalid main file")
	}
	mainContent, err := os.ReadFile(fileName)
	if err != nil {
		return &v8go.Value{}, err
	}

	e := engine.New(manifest, folder)
	return e.RunScript(string(mainContent), fileName)
}
