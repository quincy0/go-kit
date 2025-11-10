package lang

import (
	"fmt"
	"os"
	"testing"
)

func TestNewLangWithFile(t *testing.T) {
	dirPath, _ := os.Getwd()
	filePath := fmt.Sprintf("%s/lang", dirPath)

	_ = NewLangWithFile(filePath)
}
