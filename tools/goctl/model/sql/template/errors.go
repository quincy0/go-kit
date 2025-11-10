package template

// Error defines an error template
const Error = `package {{.pkg}}

import "go-kit/core/stores/sqlx"

var ErrNotFound = sqlx.ErrNotFound
`
