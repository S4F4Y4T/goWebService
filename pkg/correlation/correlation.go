package correlation

import "context"

type contextKey string

const CorrelationIDKey contextKey = "correlation_id"

// GetCorrelationID extracts the correlation ID from the context.
// This lives in pkg/ so it can be used by both internal/ and config/ without creating a boundary violation.
func GetCorrelationID(ctx context.Context) string {
	if id, ok := ctx.Value(CorrelationIDKey).(string); ok {
		return id
	}
	return ""
}
