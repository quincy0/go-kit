package lang

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

type Lang struct {
	lang      map[string]map[string]string
	langRange []string
}

func NewLangWithJson(jsonData string) *Lang {
	l := &Lang{}

	if err := json.Unmarshal([]byte(jsonData), l); err != nil {
		log.Fatal(err)
	}

	return l
}

//NewLangWithFile
//@Description: 通过文件解析多语言json文件
//@param filePath 文件路径
//@return *Lang
func NewLangWithFile(filePath string) *Lang {
	l := &Lang{}
	dir, err := ioutil.ReadDir(filePath)
	if err != nil {
		log.Fatal(err)
	}

	languageData := make(map[string]map[string]string)
	var langRange []string
	for _, file := range dir {
		mapData := make(map[string]string)
		nameArr := strings.Split(file.Name(), ".")
		if len(nameArr) >= 2 && nameArr[1] == "json" {
			languageArr := strings.Split(nameArr[0], "_")
			if len(languageArr) >= 2 {
				fileData, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", filePath, file.Name()))
				if err != nil {
					log.Fatal(err)
				}
				err = json.Unmarshal(fileData, &mapData)
				if err != nil {
					log.Fatal(err)
				}
				languageData[languageArr[1]] = mapData
				langRange = append(langRange, languageArr[1])
			}
		}
	}

	l.lang = languageData
	l.langRange = langRange

	return l
}

func (l *Lang) ParseMsg(lang, code string, v ...interface{}) string {
	if len(lang) == 0 || !inString(lang, l.langRange) {
		lang = "en"
	}
	if msgMap, ok := l.lang[lang]; ok {
		if msg, ok := msgMap[code]; ok {
			return fmt.Sprintf(msg, v...)
		}
	}
	return ""
}

func inString(target string, strArray []string) bool {
	for _, element := range strArray {
		if target == element {
			return true
		}
	}
	return false
}
