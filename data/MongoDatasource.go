package data

import (
	"context"
	"fmt"
	"github.com/hsedjame/webfast/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//MongoDatasource represents an MongoDB datasource
type MongoDatasource struct {
	WithProperties
}

//Start see Datasource.Start
func (m MongoDatasource) Start(ctx context.Context) (context.Context, error) {

	if m.IsValid() {

		// Create mongo connect uri
		uri := fmt.Sprintf("mongodb://%s:%s@%s:%d/%s",
			m.Properties.Username, m.Properties.Password, m.Properties.Host, m.Properties.Port, m.Properties.Database)

		// Create client options
		clientOptions := options.Client().ApplyURI(uri)

		if client, err := mongo.Connect(ctx, clientOptions); err != nil {
			return ctx, err
		} else {

			if err := client.Ping(ctx, nil); err != nil {
				return ctx, err
			}

			// Update the application context
			ctx = context.WithValue(ctx, utils.DBCtxKey, client)

			return ctx, nil
		}

	} else {
		return ctx, utils.AppError{
			Message: "Database is not set. Please consider to set property { \"db\" : {\"database\" : ****}}",
		}
	}
}

//Stop see Datasource.Stop
func (m MongoDatasource) Stop(ctx context.Context) {
	if client := LoadDB[*mongo.Client](ctx); client != nil {
		_ = client.Disconnect(ctx)
	}
}

//CanStart see Datasource.CanStart
func (m MongoDatasource) CanStart() bool {
	return m.IsValid()
}

//CreateTables see Datasource.CreateTables
func (m MongoDatasource) CreateTables(ctx context.Context) error {

	client := ctx.Value(utils.DBCtxKey).(*mongo.Client)

	for _, entity := range ctx.Value(utils.EntitiesCtxKey).([]Entity) {

		if err := client.Database(m.Properties.Database).CreateCollection(ctx, entity.TableName(), &options.CreateCollectionOptions{}); err != nil {
			return err
		}
	}

	return nil
}
