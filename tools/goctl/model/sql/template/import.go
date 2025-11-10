package template

const (
	// Imports defines a import template for model in cache case
	Imports = `import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	{{if .time}}"time"{{end}}

	"github.com/quincy0/go-kit/core/stores/builder"
	"github.com/quincy0/go-kit/core/stores/cache"
	"github.com/quincy0/go-kit/core/stores/sqlc"
	"github.com/quincy0/go-kit/core/stores/sqlx"
	"github.com/quincy0/go-kit/core/stringx"
)
`
	// ImportsNoCache defines a import template for model in normal case
	ImportsNoCache = `import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	{{if .time}}"time"{{end}}

	"github.com/quincy0/go-kit/core/stores/builder"
	"github.com/quincy0/go-kit/core/stores/sqlc"
	"github.com/quincy0/go-kit/core/stores/sqlx"
	"github.com/quincy0/go-kit/core/stringx"
)
`
)
