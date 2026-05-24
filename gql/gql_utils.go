package gql

import (
	"context"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/satisfactorymodding/smr-api/generated"
	"github.com/satisfactorymodding/smr-api/util"
)

// SetINN sets target if value not nil
func SetINN[T any](v *T, target *T) {
	if v != nil {
		*target = *v
	}
}

// SetCompatibilityINNF calls function 'target' to set value if not nil
func SetCompatibilityINNF[B any](value *generated.CompatibilityInfoInput, target func(*util.CompatibilityInfo) B) {
	if value == nil {
		return
	}
	target(GenCompInfoToDBCompInfo(value))
}

// SetAIDisclosureINNF calls function 'target' to set value if not nil
func SetAIDisclosureINNF[B any](value *generated.AIUseDisclosureInput, target func(*util.AIUseDisclosureInfo) B) {
	if value == nil {
		return
	}
	target(GenAIDisclosureInfoToDBAIDisclosureInfo(value))
}

// SetINNF calls function 'target' to set value if not nil
func SetINNF[T any, B any](value *T, target func(T) B) {
	if value != nil {
		target(*value)
	}
}

// SetINNOEF calls function 'target' to set value if not nil or empty
func SetINNOEF[T comparable, B any](value *T, target func(T) B) {
	if value != nil && *value != *(new(T)) {
		target(*value)
	}
}

// SetStabilityINNF calls function 'target' to set value if not nil
func SetStabilityINNF[B any](value *generated.VersionStabilities, target func(util.Stability) B) {
	if value != nil {
		target(util.Stability(*value))
	}
}

// SetDateINNF calls function 'target' to set value if not nil
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
