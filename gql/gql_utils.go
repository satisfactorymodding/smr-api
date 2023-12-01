package gql

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/satisfactorymodding/smr-api/db/postgres"
	"github.com/satisfactorymodding/smr-api/generated"
	"github.com/satisfactorymodding/smr-api/generated/ent/smlversion"
	"github.com/satisfactorymodding/smr-api/util"
)

type TraceWrapper struct {
	Span trace.Span
}

func WrapQueryTrace(ctx context.Context, action string) (TraceWrapper, context.Context) {
	return wrapTrace(ctx, action, "query")
}

func WrapMutationTrace(ctx context.Context, action string) (TraceWrapper, context.Context) {
	return wrapTrace(ctx, action, "mutation")
}

func wrapTrace(ctx context.Context, action string, actionType string) (TraceWrapper, context.Context) {
	spanCtx, span := otel.Tracer("gql").Start(ctx, "GraphQL "+action, trace.WithAttributes(
		attribute.String("action_type", "API.graphql."+actionType),
	))

	return TraceWrapper{
		Span: span,
	}, spanCtx
}

func (wrapper TraceWrapper) end() {
	defer wrapper.Span.End()

	if err := recover(); err != nil {
		wrapper.Span.RecordError(fmt.Errorf("panic: %v", err))
		panic(err)
	}
}

// SetStringINNOE sets target if value not nil or empty
func SetStringINNOE(value *string, target *string) {
	if value == nil || *value == "" {
		return
	}

	*target = *value
}

// SetINN sets target if value not nil
func SetINN[T any](v *T, target *T) {
	if v != nil {
		*target = *v
	}
}

func SetStabilityINN(value *generated.VersionStabilities, target *string) {
	if value == nil {
		return
	}

	*target = string(*value)
}

func SetCompatibilityINN(value *generated.CompatibilityInfoInput, target **postgres.CompatibilityInfo) {
	if value == nil {
		return
	}
	toDB := GenCompInfoToDBCompInfo(value)
	*target = toDB
}

func SetINNF[T any, B any](value *T, target func(T) B) {
	if value != nil {
		target(*value)
	}
}

func SetINNOEF[T comparable, B any](value *T, target func(T) B) {
	if value != nil && *value != *(new(T)) {
		target(*value)
	}
}

func SetStabilityINNF[B any](value *generated.VersionStabilities, target func(smlversion.Stability) B) {
	if value != nil {
		target(smlversion.Stability(*value))
	}
}

func SetDateINNF[B any](value *string, target func(time.Time) B) {
	if value != nil {
		t, _ := time.Parse(time.RFC3339Nano, *value)
		target(t)
	}
}

func RealIP(ctx context.Context) string {
	header := ctx.Value(util.ContextHeader{}).(http.Header)

	if ip := header.Get("X-Forwarded-For"); ip != "" {
		return strings.Split(ip, ", ")[0]
	}

	if ip := header.Get("X-Real-IP"); ip != "" {
		return ip
	}

	request := ctx.Value(util.ContextRequest{}).(*http.Request)
	ra, _, _ := net.SplitHostPort(request.RemoteAddr)

	return ra
}
