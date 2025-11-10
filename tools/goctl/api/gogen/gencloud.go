package gogen

import (
	_ "embed"
)

const CLOUD_BUILD_FILE_NAME = "cloudbuild.yaml"

//go:embed cloudbuild.tpl
var cloudTemplate string

func genCloud(dir string) error {
	return genFile(fileGenConfig{
		dir:             dir,
		subdir:          "",
		filename:        CLOUD_BUILD_FILE_NAME,
		templateName:    "cloudTemplate",
		category:        category,
		templateFile:    cloudTemplateFile,
		builtinTemplate: cloudTemplate,
		data:            map[string]string{},
	})
}
