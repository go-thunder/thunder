# Cookbook

## GraphQL

### Resolvers

- Add schema to `/internal/transport-inbound/graphql/<name>.graphql`. Note
  that the name of the file is arbitrary and impacts the name of the resolver
  file.
- Run `make generate` to generate the resolver file in
  `/internal/transport-inbound/graphql/resolvers/<name>.resolvers.go`.
- Write the code for the resolver in the generated file. Usually, this uses
  `commands` and `queries` modules from `/internal/features`.

### Loaders

```go
// package/internal/features/loaders/module.go

package loaders

import "go.uber.org/fx"

var Module = fx.Provide(
	newLoader,
)
```

```go
// package/internal/features/loaders/loaders.go

package loaders

import (
	"github.com/example/package/internal/generated/ent"
	"github.com/graph-gophers/dataloader"
)

type Loader struct {
    ...
	FindEntityByIDLoader *dataloader.Loader
	...
}

func newLoader(repo *ent.Client) *Loader {
	FindEntityByIDLoader := dataloader.NewBatchedLoader(
		findEntityByIDBatch(repo),
		dataloader.WithClearCacheOnBatch(),
	)

	return &Loader{
		...
		FindEntityByIDLoader: FindEntityByIDLoader,
		...
	}
}
```

```go
// package/internal/features/loaders/findEntityByID.go

package loaders

import (
	"context"

	"github.com/TheRafaBonin/roxy"

	"github.com/example/package/internal/generated/ent"
	"github.com/example/package/internal/generated/ent/entity"

	"github.com/google/uuid"

	"github.com/graph-gophers/dataloader"
)

func findEntityByIDBatch(repo *ent.Client) dataloader.BatchFunc {
	batchFn := func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		// Declares some variables
		var entityMap = make(map[string]*ent.Entity)
		var errorsMap = make(map[string]error)

		var results []*dataloader.Result
		var entityIDs []uuid.UUID

		// Convert keys
		for _, key := range keys.Keys() {
			uid, err := uuid.Parse(key)
			if err != nil {
				errorsMap[key] = roxy.Wrap(err, "parsing uuids")
				continue
			}

			entityIDs = append(entityIDs, uid)
		}

		// Finds entities and maps
		entities, err := repo.Entity.Query().
			Where(entity.IDIn(entityIDs...)).
			All(ctx)
		if err != nil {
			return findEntityByIDBatchErrorResults(keys, roxy.Wrap(err, "finding entities"))
		}
		for _, entity := range entities {
			entityMap[entity.ID.String()] = entity
		}

		// Map the results
		for _, key := range keys.Keys() {
			err = errorsMap[key]
			if err != nil {
				results = append(results, &dataloader.Result{
					Data:  nil,
					Error: err,
				})
				continue
			}

			p := entityMap[key]
			results = append(results, &dataloader.Result{
				Data: p,
			})
		}
		return results
	}
	return batchFn
}

func findEntityByIDBatchErrorResults(keys dataloader.Keys, err error) []*dataloader.Result {
	var results []*dataloader.Result
	for range keys.Keys() {
		results = append(results, &dataloader.Result{
			Data:  nil,
			Error: err,
		})
	}

	return results
}
```

```go
// package/internal/features/graphql/resolvers/module.go
package resolvers

import (
	"github.com/99designs/gqlgen/graphql"
	"github.com/example/package/internal/features/loaders"
	"go.uber.org/fx"

	generatedGraphql "github.com/example/package/internal/generated/graphql"
)

func newSchema(..., loaders *loaders.Loader, ...) graphql.ExecutableSchema {
	resolver := &Resolver{
		...
        loaders:  loaders,
		...
	}

	graphqlConfig := generatedGraphql.Config{Resolvers: resolver}
	return generatedGraphql.NewExecutableSchema(graphqlConfig)
}

var Module = fx.Options(
	fx.Provide(
		newSchema,
	),
)
```

```go
// package/internal/features/graphql/resolvers/resolvers.go
package resolvers

import (
	"github.com/example/package/internal/features/loaders"

	generatedGraphql "github.com/example/package/internal/generated/graphql"
)

...

type Resolver struct {
	...
    loaders  *loaders.Loader
	...
}

...
```

```go
// package/internal/features/graphql/resolvers/entity.resolvers.go

package resolvers

import (
	"context"

	"github.com/example/package/internal/generated/ent"
	"github.com/example/package/internal/generated/graphql"
	"github.com/example/package/internal/transport-inbound/graphql/resolvers/formatters"
	"github.com/graph-gophers/dataloader"
	"github.com/rotisserie/eris"
)

func (r *entityResolver) FindEntityByID(ctx context.Context, id string) (*graphql.Entity, error) {
	thunk := r.loaders.FindEntityByIDLoader.Load(ctx, dataloader.StringKey(id))
	loadedEntity, err := thunk()
	err = eris.Wrap(err, "loaders.FindEntityByIDLoader")

	if err != nil {
		return nil, err
	}

	entity := loadedEntity.(*ent.Entity)
	if entity == nil {
		return nil, err
	}

	return formatters.FormatEntity(entity), nil
}

func (r *Resolver) Entity() graphql.EntityResolver { return &entityResolver{r} }

type entityResolver struct{ *Resolver }
```

## Repository Pattern

### ent

```go
// package/cmd/entc.go

//go:build exclude

package main

import (
	"log"

	"entgo.io/contrib/entgql"
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
)

func main() {
	ex, err := entgql.NewExtension()
	if err != nil {
		log.Fatalf("creating entgql extension: %v", err)
	}

	opts := []entc.Option{
		entc.Extensions(ex),
	}

	err = entc.Generate("../internal/entities", &gen.Config{
		Target:  "../internal/generated/ent",
		Schema:  "github.com/example/package/internal/entities",
		Package: "github.com/example/package/internal/generated/ent",
		Features: []gen.Feature{
			gen.FeatureVersionedMigration,
			gen.FeatureUpsert,
			gen.FeatureLock,
		},
	}, opts...)
	if err != nil {
		log.Fatalf("running ent codegen: %v", err)
	}
}
```

## Tests

### Mock subscribers

```go
import (
    ...
	"github.com/gothunder/thunder/pkg/events/mocks"
    ...
)

...

var handler *mocks.Handler

...

// inside test case
    ...

    subscriberCalled := make(chan interface{})
    topic := events.SubscriberTopic
    var subscriberReturn events.SubscriberMessage
    handler.Mock.On("Handle",
        mock.Anything,
        topic,
        mock.Anything,
    ).Once().Return(thunderEvents.Success).Run(func(args mock.Arguments) {
        defer GinkgoRecover()

        decoder := args.Get(2).(thunderEvents.EventDecoder)
        err := decoder.Decode(&subscriberReturn)
        Expect(err).ToNot(HaveOccurred())

        close(subscriberCalled)
    })

    resp := someCommand() // this publishes an event

    // Assert
    Eventually(subscriberCalled, "3s").Should(BeClosed())
    Expect(resp).To(Equal(thunderEvents.Success))

    ...
```

## Generate command

```go
// package/cmd/generate.go

package generate

//go:generate go run entc.go
//go:generate go run github.com/99designs/gqlgen
//go:generate sh -c "protoc --experimental_allow_proto3_optional --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative --proto_path=../transport-inbound/grpc/proto --go_out=../../pkg/grpc/ --go-grpc_out=../../pkg/grpc/ ../transport-inbound/grpc/proto/*.proto"
```
