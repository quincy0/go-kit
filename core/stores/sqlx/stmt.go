package sqlx

import (
	"context"
	"database/sql"
	"time"

	"github.com/quincy0/go-kit/core/logx"
	"github.com/quincy0/go-kit/core/syncx"
	"github.com/quincy0/go-kit/core/timex"
)

const defaultSlowThreshold = time.Millisecond * 500

var slowThreshold = syncx.ForAtomicDuration(defaultSlowThreshold)

var logSqlSwitch = syncx.ForAtomicBool(true)

// SetSlowThreshold sets the slow threshold.
func SetSlowThreshold(threshold time.Duration) {
	slowThreshold.Set(threshold)
}

func DisableLog() {
	logSqlSwitch.Set(false)
}

func exec(ctx context.Context, conn sessionConn, q string, args ...interface{}) (sql.Result, error) {
	stmt, err := format(q, args...)
	if err != nil {
		return nil, err
	}

	guard := newGuard("exec", stmt)
	result, err := conn.ExecContext(ctx, q, args...)
	guard.logger(ctx)
	if err != nil {
		logSqlError(ctx, stmt, err)
	}

	return result, err
}

func execStmt(ctx context.Context, conn stmtConn, q string, args ...interface{}) (sql.Result, error) {
	stmt, err := format(q, args...)
	if err != nil {
		return nil, err
	}

	guard := newGuard("execStmt", stmt)
	result, err := conn.ExecContext(ctx, args...)
	guard.logger(ctx)
	if err != nil {
		logSqlError(ctx, stmt, err)
	}

	return result, err
}

func query(ctx context.Context, conn sessionConn, scanner func(*sql.Rows) error,
	q string, args ...interface{}) error {
	stmt, err := format(q, args...)
	if err != nil {
		return err
	}

	guard := newGuard("query", stmt)
	rows, err := conn.QueryContext(ctx, q, args...)
	guard.logger(ctx)
	if err != nil {
		logSqlError(ctx, stmt, err)
		return err
	}
	defer rows.Close()

	return scanner(rows)
}

func queryStmt(ctx context.Context, conn stmtConn, scanner func(*sql.Rows) error,
	q string, args ...interface{}) error {
	stmt, err := format(q, args...)
	if err != nil {
		return err
	}

	guard := newGuard("queryStmt", stmt)
	rows, err := conn.QueryContext(ctx, args...)
	guard.logger(ctx)
	if err != nil {
		logSqlError(ctx, stmt, err)
		return err
	}
	defer rows.Close()

	return scanner(rows)
}

type (
	sqlGuard interface {
		logger(ctx context.Context)
	}

	realSqlGuard struct {
		command   string
		stmt      string
		startTime time.Duration
	}
)

func newGuard(command, stmt string) sqlGuard {
	return &realSqlGuard{
		command:   command,
		stmt:      stmt,
		startTime: timex.Now(),
	}
}

func (e *realSqlGuard) logger(ctx context.Context) {
	duration := timex.Since(e.startTime)
	if duration > slowThreshold.Load() {
		logx.WithContext(ctx).WithDuration(duration).Slowf("[SQL] %s: slowcall - %s", e.command, e.stmt)
	} else {
		if logSqlSwitch.True() {
			logx.WithContext(ctx).WithDuration(duration).Infof("sql %s: %s", e.command, e.stmt)
		}

	}
}
