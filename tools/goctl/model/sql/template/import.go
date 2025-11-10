package template

const (
	// Imports defines a import template for model in cache case
	Imports = `import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	{{if .time}}"time"{{end}}

	"go-kit/core/stores/builder"
	"go-kit/core/stores/cache"
	"go-kit/core/stores/sqlc"
	"go-kit/core/stores/sqlx"
	"go-kit/core/stringx"
)
`
	// ImportsNoCache defines a import template for model in normal case
	ImportsNoCache = `import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	{{if .time}}"time"{{end}}

	"go-kit/core/stores/builder"
	"go-kit/core/stores/sqlc"
	"go-kit/core/stores/sqlx"
	"go-kit/core/stringx"
)
`
)
