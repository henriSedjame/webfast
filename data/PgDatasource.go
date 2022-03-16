package data

import (
	"context"
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/hsedjame/webfast/utils"
)

//PgDatasource represents an PostgreSql datasource
type PgDatasource struct {
	WithProperties
}

//Start see Datasource.Start
func (p PgDatasource) Start(ctx context.Context) (context.Context, error) {

	if p.IsValid() {

		// Create pqsl url
		optStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
			p.Properties.Username, p.Properties.Password, p.Properties.Host, p.Properties.Port, p.Properties.Database)

		// Try to parse url
		if opt, err := pg.ParseURL(optStr); err != nil {
			return ctx, err
		} else {

			// Connect the database
			db := pg.Connect(opt)

			// Try to query psql version
			// in order to verify if database is working
			var version string
			if _, err := db.QueryOneContext(ctx, pg.Scan(&version), "SELECT version()"); err != nil {
				return ctx, err
			} else {

				// Update the application context
				ctx = context.WithValue(ctx, utils.DBCtxKey, db)
				return ctx, nil
			}
		}
	} else {
		return ctx, utils.AppError{
			Message: "Database is not set. Please consider to set property { \"db\" : {\"database\" : ****}}",
		}
	}
}

//Stop see Datasource.Stop
func (p PgDatasource) Stop(ctx context.Context) {
	if db := LoadDB[*pg.DB](ctx); db != nil {
		_ = db.Close()
	}
}

//CanStart see Datasource.CanStart
func (p PgDatasource) CanStart() bool {
	return p.IsValid()
}

//CreateTables see Datasource.CreateTables
func (p PgDatasource) CreateTables(ctx context.Context) error {

	db := ctx.Value(utils.DBCtxKey).(*pg.DB)

	// Try to create a table for each entities
	for _, entity := range ctx.Value(utils.EntitiesCtxKey).([]interface{}) {

		if err := db.Model(entity).CreateTable(&orm.CreateTableOptions{
			Temp:        false,
			IfNotExists: true,
		}); err != nil {
			return err
		}
	}

	return nil
}
