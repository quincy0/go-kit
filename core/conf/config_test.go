package conf

import (
	"fmt"
	"go-kit/rest"
	"io/ioutil"
	"os"
	"testing"

	"go-kit/core/fs"
	"go-kit/core/hash"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfig_notExists(t *testing.T) {
	assert.NotNil(t, Load("not_a_file", nil))
}

func TestLoadConfig_notRecogFile(t *testing.T) {
	filename, err := fs.TempFilenameWithText("hello")
	assert.Nil(t, err)
	defer os.Remove(filename)
	assert.NotNil(t, Load(filename, nil))
}

func TestConfigJson(t *testing.T) {
	tests := []string{
		".json",
		".yaml",
		".yml",
	}
	text := `{
	"a": "foo",
	"b": 1,
	"c": "${FOO}",
	"d": "abcd!@#$112"
}`
	for _, test := range tests {
		test := test
		t.Run(test, func(t *testing.T) {
			os.Setenv("FOO", "2")
			defer os.Unsetenv("FOO")
			tmpfile, err := createTempFile(test, text)
			assert.Nil(t, err)
			defer os.Remove(tmpfile)

			var val struct {
				A string `json:"a"`
				B int    `json:"b"`
				C string `json:"c"`
				D string `json:"d"`
			}
			MustLoad(tmpfile, &val)
			assert.Equal(t, "foo", val.A)
			assert.Equal(t, 1, val.B)
			assert.Equal(t, "${FOO}", val.C)
			assert.Equal(t, "abcd!@#$112", val.D)
		})
	}
}

func TestConfigToml(t *testing.T) {
	text := `a = "foo"
b = 1
c = "${FOO}"
d = "abcd!@#$112"
`
	os.Setenv("FOO", "2")
	defer os.Unsetenv("FOO")
	tmpfile, err := createTempFile(".toml", text)
	assert.Nil(t, err)
	defer os.Remove(tmpfile)

	var val struct {
		A string `json:"a"`
		B int    `json:"b"`
		C string `json:"c"`
		D string `json:"d"`
	}
	MustLoad(tmpfile, &val)
	assert.Equal(t, "foo", val.A)
	assert.Equal(t, 1, val.B)
	assert.Equal(t, "${FOO}", val.C)
	assert.Equal(t, "abcd!@#$112", val.D)
}

func TestConfigTomlEnv(t *testing.T) {
	text := `a = "foo"
b = 1
c = "${FOO}"
d = "abcd!@#112"
`
	os.Setenv("FOO", "2")
	defer os.Unsetenv("FOO")
	tmpfile, err := createTempFile(".toml", text)
	assert.Nil(t, err)
	defer os.Remove(tmpfile)

	var val struct {
		A string `json:"a"`
		B int    `json:"b"`
		C string `json:"c"`
		D string `json:"d"`
	}

	MustLoad(tmpfile, &val, UseEnv())
	assert.Equal(t, "foo", val.A)
	assert.Equal(t, 1, val.B)
	assert.Equal(t, "2", val.C)
	assert.Equal(t, "abcd!@#112", val.D)

}

func TestConfigJsonEnv(t *testing.T) {
	tests := []string{
		".json",
		".yaml",
		".yml",
	}
	text := `{
	"a": "foo",
	"b": 1,
	"c": "${FOO}",
	"d": "abcd!@#$a12 3"
}`
	for _, test := range tests {
		test := test
		t.Run(test, func(t *testing.T) {
			os.Setenv("FOO", "2")
			defer os.Unsetenv("FOO")
			tmpfile, err := createTempFile(test, text)
			assert.Nil(t, err)
			defer os.Remove(tmpfile)

			var val struct {
				A string `json:"a"`
				B int    `json:"b"`
				C string `json:"c"`
				D string `json:"d"`
			}
			MustLoad(tmpfile, &val, UseEnv())
			assert.Equal(t, "foo", val.A)
			assert.Equal(t, 1, val.B)
			assert.Equal(t, "2", val.C)
			assert.Equal(t, "abcd!@# 3", val.D)
		})
	}
}

func TestLoadFromConfigService(t *testing.T) {

	type Config struct {
		rest.RestConf
	}

	os.Setenv("ENVIRONMENT", "dev")

	c := Config{}

	err := LoadFromConfigService("account-conf", "account-api", &c)

	fmt.Println(c, " err:", err)
}

func createTempFile(ext, text string) (string, error) {
	tmpfile, err := ioutil.TempFile(os.TempDir(), hash.Md5Hex([]byte(text))+"*"+ext)
	if err != nil {
		return "", err
	}

	if err := ioutil.WriteFile(tmpfile.Name(), []byte(text), os.ModeTemporary); err != nil {
		return "", err
	}

	filename := tmpfile.Name()
	if err = tmpfile.Close(); err != nil {
		return "", err
	}

	return filename, nil
}
