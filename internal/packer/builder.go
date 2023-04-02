package packer

import (
	"errors"
	"github.com/evanw/esbuild/pkg/api"
	"os"
	"serverless/internal/executor"
	"strings"
)

func Build(folder string) (api.BuildResult, error) {
	manifest, err := executor.ParseManifest(folder)
	if err != nil {
		return api.BuildResult{}, err
	}

	outDir := folder + "/" + manifest.BuildDir
	buildResult := api.Build(api.BuildOptions{
		EntryPoints:       []string{strings.Trim(folder+"/"+manifest.Main, string(os.PathSeparator))},
		Outdir:            strings.Trim(outDir, string(os.PathSeparator)),
		Bundle:            true,
		Write:             true,
		MinifyWhitespace:  true,
		MinifyIdentifiers: true,
		MinifySyntax:      true,
		JSXFactory:        "h",
		Platform:          api.PlatformNeutral,
		LogLevel:          api.LogLevelInfo,
	})
	errorMessages := buildResult.Errors
	if len(errorMessages) > 0 {
		var messages string
		for i := range errorMessages {
			message := errorMessages[i]
			messages += message.Text + "\n"
		}
		return buildResult, errors.New(messages)
	}
	return buildResult, nil
}
