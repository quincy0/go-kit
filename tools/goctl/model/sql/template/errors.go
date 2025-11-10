package template

// Error defines an error template
const Error = `package {{.pkg}}

import "github.com/quincy0/go-kit/core/stores/sqlx"

var ErrNotFound = sqlx.ErrNotFound
`
