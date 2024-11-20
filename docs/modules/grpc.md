# GRPC

## Prerequisites

Install the prerequisites listed in [gRPC Go - Quick Start](https://grpc.io/docs/languages/go/quickstart/#prerequisites).

## Getting Started

Add the following code to `cmd/generate/generate.go`:

```
//go:generate sh -c "protoc --experimental_allow_proto3_optional --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative --proto_path=../transport-inbound/grpc/proto --go_out=../../pkg/grpc/ --go-grpc_out=../../pkg/grpc/ ../transport-inbound/grpc/proto/*.proto"
```

In `main.go`:

```go
    ...
	fx.New(
		fx.Populate(&w),

		// Adapters
		thunderRabbitMQ.PublisherModule,
		commonsConnection.Native,
		thunderRouter.Module,
		thunderLog.Module,
		thunderGraphql.Module,

		// Internal Modules
		transportoutbound.Module,
		transportinbound.Module,
		generated.Module,
		features.Module,

		// Listeners
		fx.Invoke(
			commonsMigrate.MigrateUp,
			auth.RegisterMiddleware,
			thunderRouter.StartListener,
			thunderGrpc.StartGrpcListener,   <----
		),
	).Run()
    ...
```

Inside `internal/transport-inbound/module.go`, add the gRPC module:

```go
package transportinbound

import (
	"github.com/example/package/internal/transport-inbound/grpc"

	"go.uber.org/fx"
)

var Module = fx.Options(
	...
	grpc.Module,
    ...
)
```

Inside the gRPC module, provide the gRPC server:

```go
// internal/transport-inbound/grpc/server.go
package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/example/package/v3/internal/features/commands"
	"github.com/example/package/internal/features/queries"
	pb "github.com/example/package/pkg/grpc"
	featureserver "github.com/example/package/internal/transport-inbound/grpc/feature-server"
	"github.com/rs/zerolog"

	thunderGrpc "github.com/gothunder/thunder/pkg/grpc"

	"google.golang.org/grpc"
)

type server struct {
	FeatureServer featureserver.FeatureServer
	grpcServer    *grpc.Server
}

// NewGrpcServer creates a new grpc server
func NewGrpcServer(commands *commands.CommandGroup, queries *queries.QueryGroup, logger *zerolog.Logger) thunderGrpc.GrpcServer {
	// Declares the server
	server := instantiateNewGrpcServer(commands, queries, logger)

	// Register service servers
	pb.RegisterFeatureServiceServer(server.grpcServer, server.FeatureServer)

	return server
}

func instantiateNewGrpcServer(commands *commands.CommandGroup, queries *queries.QueryGroup, logger *zerolog.Logger) server {
	server := server{
		FeatureServer: featureserver.NewFeatureServer(commands, queries),
	}

	sv := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpcLoggerInterceptor(logger),
		),
	)
	server.grpcServer = sv

	return server
}

func (s server) GetListener() (net.Listener, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 50051))
	if err != nil {
		return nil, err
	}

	return lis, nil
}

func (s server) GetGrpcServer() *grpc.Server {
	return s.grpcServer
}

func grpcLoggerInterceptor(logger *zerolog.Logger) func(context.Context, interface{}, *grpc.UnaryServerInfo, grpc.UnaryHandler) (interface{}, error) {

	return func(ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {

		// Add logger to context
		lctx := logger.WithContext(ctx)

		// Calls the handler
		h, err := handler(lctx, req)

		return h, err
	}
}
```

In this case, we are using the `FeatureServer` from the `featureserver`
package. This struct implements the `FeatureServiceServer` interface generated
by the protobuf compiler from a `feature.proto` file.

`commands` and `queries` can be omitted or replaced by any other dependency.

As for the `feature-server` package, it should look like this:

```go
// internal/transport-inbound/grpc/feature-server/server.go

package featureserver

import (
	"github.com/example/package/internal/features/commands"
	"github.com/example/package/internal/features/queries"
	pb "github.com/example/package/pkg/grpc"
)

// FeatureServer defines the feature server
type FeatureServer struct {
	pb.UnimplementedFeatureServiceServer
	commands *commands.CommandGroup
	queries  *queries.QueryGroup
}

// NewFeatureServer creates a new grpc server
func NewFeatureServer(commands *commands.CommandGroup, queries *queries.QueryGroup) FeatureServer {
	server := FeatureServer{
		commands: commands,
		queries:  queries,
	}

	return server
}
```

```go
// internal/transport-inbound/grpc/feature-server/rpcName.go

package featureserver

import (
	"context"

	"github.com/TheRafaBonin/roxy"
	thunderGrpc "github.com/gothunder/thunder/pkg/grpc"

	"github.com/example/package/internal/errors"
	"github.com/example/package/internal/features/queries/filters"
	"github.com/example/package/internal/features/queries/relations"
	"github.com/example/package/internal/transport-inbound/grpc/feature-server/formatters"
	"github.com/example/package/internal/transport-inbound/grpc/feature-server/parsers"

	pb "github.com/example/package/pkg/grpc"
)

// ListFeatureEntities ...
func (s FeatureServer) ListFeatureEntities(ctx context.Context, in *pb.ListFeatureEntitiesRequest) (*pb.ListFeatureEntitiesReply, error) {
	// Validates input
	if in == nil {
		statusErr := thunderGrpc.HandleError(
        ctx, 
        roxy.SetDefaultGrpcResponse(roxy.New("Input cannot be nil"), roxy.GrpcResponse{
            Message: "Invalid input provided",
            Code:    codes.InvalidArgument,
        })

		return nil, statusErr
	}

	// Parse variables
	parsedInput, err := parsers.ParseEntityFilters(in.EntityFilters)
	err = roxy.Wrap(err, "parsers.ParseEntityFilters")
	if err != nil {
		statusErr := thunderGrpc.HandleError(ctx, err)
		return nil, statusErr
	}

	// Gets entities
	entities, err := s.queries.Entity.FindMany(
        ctx, 
        &filters.EntityFilterInput{ ... },
    )
	err = roxy.Wrap(err, "queries.Entity.FindMany")
	if err != nil {
		statusErr := thunderGrpc.HandleError(ctx, err)
		return nil, statusErr
	}

	// Formats response
	formattedEntities, err := formatters.FormatEntities(ctx, Entities)
	err = roxy.Wrap(err, "formatters.FormatEntities")
	if err != nil {
		statusErr := thunderGrpc.HandleError(ctx, err)
		return nil, statusErr
	}

	return &pb.ListFeatureEntitiesReply{
		Entities: formattedEntities,
	}, nil
}
```
