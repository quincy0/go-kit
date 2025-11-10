package gogen

import (
	_ "embed"

	"github.com/quincy0/go-kit/tools/goctl/api/spec"
)

const DOCKERFILE_FILE_NAME = "Dockerfile"

//go:embed dockerfile.tpl
var dockerfileTemplate string

func genDockerfile(dir string, api *spec.ApiSpec) error {
	return genFile(fileGenConfig{
		dir:             dir,
		subdir:          "",
		filename:        DOCKERFILE_FILE_NAME,
		templateName:    "dockerfileTemplate",
		category:        category,
		templateFile:    dockerfileTemplateFile,
		builtinTemplate: dockerfileTemplate,
		data: map[string]string{
			"APP_NAME": api.Service.Name,
		},
	})
}
