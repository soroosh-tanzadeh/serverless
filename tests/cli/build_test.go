package cli

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
	"os"
	"serverless/cli/platform"
	"testing"
)

var app *cli.App

func TestMain(m *testing.M) {
	app = platform.CreateApp()
	m.Run()
}

func TestBuildShouldBuildJsFileCorrectly(t *testing.T) {
	err := app.Command("build").Run(cli.NewContext(app, nil, &cli.Context{Context: context.Background()}), "-f", "../test-app")
	if err != nil {
		t.Error(err)
	}
	var buildPath = "../test-app/build/app.js"
	assert.FileExists(t, buildPath)
	var expectedContent = "response(`<h1>Hello World</h1><button onclick=\"alert('Hello')\">Click Me</button>`,200,{\"content-type\":\"text/html\",\"x-cache\":!0,\"x-server\":1025});\n"
	actualContent, err := os.ReadFile(buildPath)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, expectedContent, string(actualContent))
}
