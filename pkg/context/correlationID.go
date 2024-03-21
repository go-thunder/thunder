package context

import (
	"context"

	"github.com/google/uuid"
)

const ThunderCorrelationIDMetadataKey = "x-thunder-correlation-id"

// CorrelationIDFromContext retrieves the correlation ID from a context.Context.
func CorrelationIDFromContext(ctx context.Context) string {
	md := MetadataFromContext(ctx)
	if md == nil {
		// creates a new correlation ID if none is found.
		// it makes sure correlation is passed through the system even if it was not set
		// by the first caller
		return uuid.Must(uuid.NewV7()).String()
	}

	return md.Get(ThunderCorrelationIDMetadataKey)
}

// ContextWithCorrelationID returns a new context.Context that holds the given correlation ID.
func ContextWithCorrelationID(ctx context.Context, correlationID string) context.Context {
	md := make(Metadata, 1)
	if correlationID == "" {
		md.Set(ThunderCorrelationIDMetadataKey, uuid.Must(uuid.NewV7()).String())
	} else {
		md.Set(ThunderCorrelationIDMetadataKey, correlationID)
	}
	return ContextWithMetadata(ctx, md)
}
