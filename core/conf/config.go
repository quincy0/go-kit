package conf

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"go-kit/core/mapping"
	"go-kit/rest/httpc"
)

var loaders = map[string]func([]byte, interface{}) error{
	".json": LoadFromJsonBytes,
	".toml": LoadFromTomlBytes,
	".yaml": LoadFromYamlBytes,
	".yml":  LoadFromYamlBytes,
}

type ConfItem struct {
	ID            int    `json:"id"`
	AppName       string `json:"appName"`
	ClusterName   string `json:"clusterName"`
	NamespaceName string `json:"namespaceName"`
	Format        string `json:"format"`
	Value         string `json:"value"`
	CreatedAt     int    `json:"createdAt"`
	UpdatedAt     int    `json:"updatedAt"`
}
type ConfReply struct {
	Code int      `json:"code"`
	Msg  string   `json:"msg"`
	Data ConfItem `json:"data"`
}

// Load loads config into v from file, .json, .yaml and .yml are acceptable.
func Load(file string, v interface{}, opts ...Option) error {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	loader, ok := loaders[strings.ToLower(path.Ext(file))]
	if !ok {
		return fmt.Errorf("unrecognized file type: %s", file)
	}

	var opt options
	for _, o := range opts {
		o(&opt)
	}

	if opt.env {
		return loader([]byte(os.ExpandEnv(string(content))), v)
	}

	return loader(content, v)
}

// LoadConfig loads config into v from file, .json, .yaml and .yml are acceptable.
// Deprecated: use Load instead.
func LoadConfig(file string, v interface{}, opts ...Option) error {
	return Load(file, v, opts...)
}

// LoadFromJsonBytes loads config into v from content json bytes.
func LoadFromJsonBytes(content []byte, v interface{}) error {
	return mapping.UnmarshalJsonBytes(content, v)
}

// LoadConfigFromJsonBytes loads config into v from content json bytes.
// Deprecated: use LoadFromJsonBytes instead.
func LoadConfigFromJsonBytes(content []byte, v interface{}) error {
	return LoadFromJsonBytes(content, v)
}

// LoadFromTomlBytes loads config into v from content toml bytes.
func LoadFromTomlBytes(content []byte, v interface{}) error {
	return mapping.UnmarshalTomlBytes(content, v)
}

// LoadFromYamlBytes loads config into v from content yaml bytes.
func LoadFromYamlBytes(content []byte, v interface{}) error {
	return mapping.UnmarshalYamlBytes(content, v)
}

// LoadConfigFromYamlBytes loads config into v from content yaml bytes.
// Deprecated: use LoadFromYamlBytes instead.
func LoadConfigFromYamlBytes(content []byte, v interface{}) error {
	return LoadFromYamlBytes(content, v)
}

// MustLoad loads config into v from path, exits on error.
func MustLoad(path string, v interface{}, opts ...Option) {
	if err := Load(path, v, opts...); err != nil {
		log.Fatalf("error: config file %s, %s", path, err.Error())
	}
}

func LoadFromConfigService(clusterName string, serviceName string, v interface{}) error {
	env := os.Getenv("ENVIRONMENT")
	projectId := os.Getenv("PROJECT_ID")
	if env != "dev" && env != "production" {
		return errors.New("ENVIRONMENT not set")
	}

	url := fmt.Sprintf("http://dev.internal.xconf/config/%s/dev/%s", clusterName, serviceName)
	if env == "production" {
		url = fmt.Sprintf("http://prod.internal.xconf/config/%s/production/%s", clusterName, serviceName)
	}
	if projectId == "gexhub" {
		url = fmt.Sprintf("https://xconf-223026273590.northamerica-northeast1.run.app/config/%s/%s/%s", clusterName, env, serviceName)
	}
	
	ret, err := httpc.Do(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	if ret != nil {
		defer ret.Body.Close()
	}
	confReply := ConfReply{}
	bodyStr, err := ioutil.ReadAll(ret.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(bodyStr, &confReply); err != nil {
		return err
	}

	if confReply.Code != 0 {
		return errors.New(fmt.Sprintf("config reply err: %s", confReply.Msg))
	}

	if env == "dev" {
		fmt.Printf("Config load result url:%s, replay:%s", url, string(bodyStr))
	}

	switch confReply.Data.Format {
	case "yaml":
		err := LoadFromYamlBytes([]byte(confReply.Data.Value), v)
		if err != nil {
			return err
		}
	case "json":
		err := LoadFromJsonBytes([]byte(confReply.Data.Value), v)
		if err != nil {
			return err
		}
	default:
		return errors.New(fmt.Sprintf("config data format invalid:%s", confReply.Data.Format))
	}

	return nil
}
