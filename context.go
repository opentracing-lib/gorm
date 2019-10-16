package gorm

import (
	apmsql "github.com/opentracing-contrib/sql"
	"context"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

const (
	apmContextKey = "apm:context"
)

// WithContext returns a copy of db with ctx recorded for use by
// the callbacks registered via RegisterCallbacks.
func WithContext(ctx context.Context, db *gorm.DB) *gorm.DB {
	return db.Set(apmContextKey, ctx)
}

func scopeContext(scope *gorm.Scope) (context.Context, bool) {
	value, ok := scope.Get(apmContextKey)
	if !ok {
		return nil, false
	}
	ctx, _ := value.(context.Context)
	return ctx, ctx != nil
}

// RegisterCallbacks registers callbacks on db for reporting spans
// to Elastic APM. This is called automatically by apmgorm.Open;
// it is provided for cases where a *gorm.DB is acquired by other
// means.
func RegisterCallbacks(db *gorm.DB) {
	registerCallbacks(db, apmsql.DSNInfo{})
}

func registerCallbacks(db *gorm.DB, dsnInfo apmsql.DSNInfo) {
	driverName := db.Dialect().GetName()
	switch driverName {
	case "postgres":
		driverName = "postgresql"
	}
	spanTypePrefix := fmt.Sprintf("db.%s.", driverName)
	querySpanType := spanTypePrefix + "query"
	execSpanType := spanTypePrefix + "exec"

	type params struct {
		spanType  string
		processor func() *gorm.CallbackProcessor
	}
	callbacks := map[string]params{
		"gorm:create": {
			spanType:  execSpanType,
			processor: func() *gorm.CallbackProcessor { return db.Callback().Create() },
		},
		"gorm:delete": {
			spanType:  execSpanType,
			processor: func() *gorm.CallbackProcessor { return db.Callback().Delete() },
		},
		"gorm:query": {
			spanType:  querySpanType,
			processor: func() *gorm.CallbackProcessor { return db.Callback().Query() },
		},
		"gorm:update": {
			spanType:  execSpanType,
			processor: func() *gorm.CallbackProcessor { return db.Callback().Update() },
		},
		"gorm:row_query": {
			spanType:  querySpanType,
			processor: func() *gorm.CallbackProcessor { return db.Callback().RowQuery() },
		},
	}
	for name, params := range callbacks {
		const callbackPrefix = "apm"
		params.processor().Before(name).Register(
			fmt.Sprintf("%s:before:%s", callbackPrefix, name),
			newBeforeCallback(params.spanType),
		)
		params.processor().After(name).Register(
			fmt.Sprintf("%s:after:%s", callbackPrefix, name),
			newAfterCallback(dsnInfo),
		)
	}
}

func newBeforeCallback(spanType string) func(*gorm.Scope) {
	return func(scope *gorm.Scope) {
		ctx, ok := scopeContext(scope)
		if !ok {
			return
		}

		// 开始span
		_, ctx = opentracing.StartSpanFromContext(ctx, spanType)

		scope.Set(apmContextKey, ctx)
	}
}

func newAfterCallback(dsnInfo apmsql.DSNInfo) func(*gorm.Scope) {
	return func(scope *gorm.Scope) {
		ctx, ok := scopeContext(scope)
		if !ok {
			return
		}

		// 完成span
		span := opentracing.SpanFromContext(ctx)
		if span == nil {
			return
		}

		ext.DBUser.Set(span, dsnInfo.User)
		ext.DBInstance.Set(span, dsnInfo.Database)
		ext.DBType.Set(span, "sql")
		ext.DBStatement.Set(span, scope.SQL)

		defer span.Finish()
	}
}
