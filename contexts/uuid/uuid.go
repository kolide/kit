// Package uuid provides helpers for propagating UUIDs accross services.
package uuid

import (
	"context"

	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
)

// Use a private type to prevent name collisions with other packages.
type key string

const uuidKey key = "UUID"

// NewContext creates a new context with the UUID set to the provided value
func NewContext(ctx context.Context, uuid string) context.Context {
	return context.WithValue(ctx, uuidKey, uuid)
}

// FromContext returns the UUID value stored in ctx, if any.
func FromContext(ctx context.Context) (string, bool) {
	uuid, ok := ctx.Value(uuidKey).(string)
	return uuid, ok
}

// NewForRequest returns a Random (Version 4) UUID string.
// If the UUID fails to generate an empty string will be returned.
func NewForRequest() string {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return ""
	}
	return uuid.String()
}

// Attach adds the UUID values stored in context to the gRPC request metadata.
func Attach() grpctransport.ClientOption {
	return grpctransport.ClientBefore(
		func(ctx context.Context, md *metadata.MD) context.Context {
			uuid, _ := FromContext(ctx)
			return grpctransport.SetRequestHeader("uuid", uuid)(ctx, md)
		},
	)
}
