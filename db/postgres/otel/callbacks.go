package otel

import (
	"strings"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	oteltrace "go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

const (
	spanName = "gorm.query"

	dbTableKey     = attribute.Key("db.table")
	dbCountKey     = attribute.Key("db.count")
	dbOperationKey = semconv.DBOperationKey
	dbStatementKey = semconv.DBStatementKey
)

func dbTable(name string) attribute.KeyValue {
	return dbTableKey.String(name)
}

func dbStatement(stmt string) attribute.KeyValue {
	return dbStatementKey.String(stmt)
}

func dbCount(n int64) attribute.KeyValue {
	return dbCountKey.Int64(n)
}

func dbOperation(op string) attribute.KeyValue {
	return dbOperationKey.String(op)
}

func (op *Plugin) before(tx *gorm.DB) {
	tx.Statement.Context, _ = op.tracer.
		Start(tx.Statement.Context, spanName, oteltrace.WithSpanKind(oteltrace.SpanKindClient))
}

func extractQuery(tx *gorm.DB) string {
	return tx.Dialector.Explain(tx.Statement.SQL.String(), tx.Statement.Vars...)
}

func (op *Plugin) after(operation string) gormHookFunc {
	return func(tx *gorm.DB) {
		span := oteltrace.SpanFromContext(tx.Statement.Context)
		if !span.IsRecording() {
			// skip the reporting if not recording
			return
		}
		defer span.End()

		// Error
		if tx.Error != nil {
			span.SetStatus(codes.Error, tx.Error.Error())
		}

		// extract the db operation
		query := extractQuery(tx)
		if operation == "" {
			operation = strings.ToUpper(strings.Split(query, " ")[0])
		}

		if tx.Statement.Table != "" {
			span.SetAttributes(dbTable(tx.Statement.Table))
		}

		span.SetAttributes(
			dbStatement(query),
			dbOperation(operation),
			dbCount(tx.Statement.RowsAffected),
		)
	}
}
