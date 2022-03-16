package data

import (
	"context"
	"fmt"
	"github.com/edgedb/edgedb-go"
	"github.com/hsedjame/webfast/utils"
)

//EdgeDatasource represents an EdgeDB datasource
type EdgeDatasource struct {
	WithProperties
}

//Start see Datasource.Start
func (e EdgeDatasource) Start(ctx context.Context) (context.Context, error) {

	if e.IsValid() {
		options := edgedb.Options{
			Host:        e.Properties.Host,
			Port:        int(e.Properties.Port),
			Database:    e.Properties.Database,
			User:        e.Properties.Username,
			Password:    edgedb.NewOptionalStr(e.Properties.Password),
			Concurrency: 4,
		}
		if client, err := edgedb.CreateClient(ctx, options); err != nil {
			return ctx, err
		} else {
			ctx = context.WithValue(ctx, utils.DBCtxKey, *client)
			return ctx, nil
		}

	} else {
		return ctx, utils.AppError{
			Message: "Database is not set. Please consider to set property { \"db\" : {\"database\" : ****}}",
		}
	}
}

//Stop see Datasource.Stop
func (e EdgeDatasource) Stop(ctx context.Context) {
	if client := LoadDB[*edgedb.Client](ctx); client != nil {
		_ = client.Close()
	}
}

//CanStart see Datasource.CanStart
func (e EdgeDatasource) CanStart() bool {
	return e.IsValid()
}

//CreateTables see Datasource.CreateTables
func (e EdgeDatasource) CreateTables(ctx context.Context) error {

	fmt.Println("Creating entities table is not supported yet with EDGEDB")
	fmt.Println("Remember to edit your schema.esdl in dbschema folder to your entities definitions")
	fmt.Println("Then run `edgedb  migration create` and `edgedb migrate` to create your tables")

	return nil
}
