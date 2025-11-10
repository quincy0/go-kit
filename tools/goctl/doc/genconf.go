package doc

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

const CONF_PATH = "~/.docrc/"

type YapiAuthConf struct {
	Type   string `json:"type"`
	Token  string `json:"token"`
	File   string `json:"file"`
	Merge  string `json:"merge"`
	Server string `json:"server"`
}

type DocConf struct {
	Server string                 `json:"server"`
	Apps   map[string]DocConfItem `json:"apps"`
}

type DocConfItem struct {
	Type   string `json:"type"`
	Token  string `json:"token"`
	AppDir string `json:"appDir"`
	Merge  string `json:"merge"`
	Server string `json:"server"`
}

func newDocConf() *DocConf {
	fd, err := os.Open(getDocrcDir() + "app.json")
	if err != nil {
		log.Fatalf("File : .docrc/app.json is not found.")
	}
	defer fd.Close()

	ret, err := os.ReadFile(getDocrcDir() + "app.json")
	if err != nil {
		log.Fatalf("read file app.json failed.")
	}

	dC := &DocConf{}

	if err := json.Unmarshal(ret, dC); err != nil {
		log.Fatalf("Unmarshal conf faile.")
	}

	return dC
}

func (d *DocConf) genDocConf() error {
	if _, err := os.Stat(getDocrcDir() + "apps/"); os.IsNotExist(err) {
		os.Mkdir(getDocrcDir()+"apps/", os.ModePerm)
	}

	for k, v := range d.Apps {
		yConf := YapiAuthConf{
			Type:   v.Type,
			Token:  v.Token,
			File:   fmt.Sprintf("%s/%s.json", v.AppDir, k),
			Merge:  v.Merge,
			Server: d.Server,
		}
		yBytes, err := json.Marshal(yConf)
		if err != nil {
			return err
		}

		if err := os.WriteFile(fmt.Sprintf("%s%s/%s.json", getDocrcDir(), "apps", k), yBytes, os.ModePerm); err != nil {
			return err
		}
	}

	return nil
}

func (d *DocConf) getDocConfItemByAppName(appName string) *DocConfItem {
	item, ok := d.Apps[appName]
	if !ok {
		return nil
	}

	return &item
}

func getDocrcDir() string {
	path, err := os.UserHomeDir()
	if err != nil || len(path) == 0 {
		path = CONF_PATH
	} else {
		path = path + "/.docrc/"
	}
	return path
}
