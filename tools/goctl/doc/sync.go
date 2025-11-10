package doc

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// DartCommand create dart network request code
func Sync(_ *cobra.Command, _ []string) error {
	genC := newDocConf()
	genC.genDocConf()

	docrcDir := getDocrcDir()

	appsDir := fmt.Sprintf("%s%s", docrcDir, "apps/")

	appConfFile := fmt.Sprintf("%s%s.json", appsDir, app)

	if len(api) == 0 {
		return errors.New("please input .api file path, like 'app.api")
	}

	if _, err := os.Stat(api); os.IsNotExist(err) {
		return errors.New(".api file is not found.")
	}

	if len(app) == 0 {
		return errors.New("the -app usage is required")
	}

	if _, err := os.Stat(appConfFile); os.IsNotExist(err) {
		return errors.New(fmt.Sprintf("app: %s is not found in ~/.docrc/app.json, please checked it.", app))
	}

	appItem := genC.getDocConfItemByAppName(app)
	if appItem == nil {
		return errors.New("app item is not found.")
	}

	// 通过goctl重新生成 swagger
	genSwaggerParmas := []string{
		"api",
		"plugin",
		"-plugin",
		fmt.Sprintf("goctl-swagger=\"swagger -filename %s.json\"", app),
		"-api",
		api,
		"-dir",
		appItem.AppDir,
	}
	cmd := exec.Command("goctl", genSwaggerParmas...)
	fmt.Println("gen swagger...\n", cmd.String())
	ret, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	fmt.Println(string(ret))

	cmd = exec.Command("yapi", "import", "--config", appConfFile)
	fmt.Println("command: ", cmd.String())
	ret, err = cmd.CombinedOutput()
	fmt.Println(string(ret))
	if err != nil {
		return err
	}

	return nil
}
