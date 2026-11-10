package middleware

import (
	"context"
	"strings"
	"sync"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"

	common "github.com/pydio/cells-runtime"
	"github.com/pydio/cells-runtime/errors"
	"github.com/pydio/cells-runtime/telemetry/log"
	"github.com/pydio/cells-runtime/utils/uuid"
)

const (
	errCauseIdKey       = "ErrorCauseId"
	errCauseMsgKey      = "ErrorCauseMsg"
	errHandlerSpanIDKey = "ErrorHandlerSpanID"
)

var errorPolicy = struct {
	sync.RWMutex
	ignores []error
	warns   []error
}{ignores: []error{errors.StatusCancelled, context.Canceled, errors.LogDAOResolver}, warns: []error{errors.HealthCheckError, errors.StatusServiceUnavailable}}

// RegisterErrorLogPolicy adds application-owned sentinels that should be
// ignored or logged as warnings by the generic gRPC error middleware.
func RegisterErrorLogPolicy(ignores, warns []error) {
	errorPolicy.Lock()
	defer errorPolicy.Unlock()
	errorPolicy.ignores = append(errorPolicy.ignores, ignores...)
	errorPolicy.warns = append(errorPolicy.warns, warns...)
}

func handleGrpcError(ctx context.Context, err error, prefix string, infos ...zap.Field) error {
	if err == nil {
		return nil
	}
	errorPolicy.RLock()
	ignores := append([]error(nil), errorPolicy.ignores...)
	warns := append([]error(nil), errorPolicy.warns...)
	errorPolicy.RUnlock()
	if len(metadata.ValueFromIncomingContext(ctx, common.CtxGrpcSilentNotFound)) > 0 {
		ignores = append(ignores, errors.StatusNotFound)
	}
	if errors.Is(err, HandledError) {
		return err
	}

	errorID := uuid.New()[:13]
	errorMessage := strings.ReplaceAll(err.Error(), "\n", " || ")
	ignored, warning := matchesAny(err, ignores), matchesAny(err, warns)
	if !ignored {
		fields := append([]zap.Field{zap.String("errorId", errorID)}, infos...)
		if caller := metadata.ValueFromIncomingContext(ctx, common.CtxGrpcClientCaller); len(caller) > 0 {
			fields = append(fields, zap.String("ClientCaller", strings.Join(caller, "")))
		}
		if warning {
			log.Logger(ctx).Warn(prefix+" "+errorMessage, fields...)
		} else {
			log.Logger(ctx).Error(prefix+" "+errorMessage, append(fields, errors.Zap(err))...)
		}
	}
	if errors.Is(err, context.Canceled) {
		err = errors.Tag(err, errors.StatusCancelled)
	}
	err = errors.Tag(err, HandledError)
	details := []any{errCauseIdKey, errorID, errCauseMsgKey, errorMessage}
	if span := trace.SpanFromContext(ctx).SpanContext(); span.IsValid() {
		details = append(details, errHandlerSpanIDKey, span.SpanID().String())
	}
	return errors.WithDetails(err, details...)
}

func matchesAny(err error, candidates []error) bool {
	for _, candidate := range candidates {
		if errors.Is(err, candidate) {
			return true
		}
	}
	return false
}
