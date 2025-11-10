module go-kit/tools/goctl

go 1.16

require (
	github.com/emicklei/proto v1.10.0
	github.com/fatih/structtag v1.2.0
	github.com/go-sql-driver/mysql v1.7.0
	github.com/google/go-github/v47 v47.1.0
	go-kit v1.11.16
	github.com/iancoleman/strcase v0.2.0
	github.com/logrusorgru/aurora v2.0.3+incompatible
	github.com/spf13/cobra v1.4.0
	github.com/stretchr/testify v1.8.2
	github.com/withfig/autocomplete-tools/integrations/cobra v0.0.0-20220705165518-2761d7f4b8bc
	github.com/zeromicro/antlr v0.0.1
	github.com/zeromicro/ddl-parser v1.0.4
	golang.org/x/oauth2 v0.4.0
	google.golang.org/grpc v1.53.0
	google.golang.org/protobuf v1.28.2-0.20220831092852-f930b1dc76e8
)

replace (
	go-kit => ../..
	go-kit/tools/goctl => ../../tools/goctl
)
