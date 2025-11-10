package gen

import (
	"github.com/quincy0/go-kit/tools/goctl/model/sql/template"
	"github.com/quincy0/go-kit/tools/goctl/util"
	"github.com/quincy0/go-kit/tools/goctl/util/pathx"
)

func genTableName(table Table) (string, error) {
	text, err := pathx.LoadTemplate(category, tableNameTemplateFile, template.TableName)
	if err != nil {
		return "", err
	}

	output, err := util.With("tableName").
		Parse(text).
		Execute(map[string]interface{}{
			"tableName":             table.Name.Source(),
			"upperStartCamelObject": table.Name.ToCamel(),
		})
	if err != nil {
		return "", nil
	}

	return output.String(), nil
}
