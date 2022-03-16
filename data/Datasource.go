package data

import (
	"context"
	"errors"
	"fmt"
)

// Datasource represents a datasource that can be started and stopped
// It can also use to create entities tables
type Datasource interface {
	//Start to start the datasource
	Start(context.Context) (context.Context, error)

	//Stop to close the datasource
	Stop(context.Context)

	//CanStart is called before Start to check if the datasource is ready to start
	CanStart() bool

	//CreateTables create all entities table
	// except for EGDEDB datasource type
	CreateTables(context.Context) error
}

// DatasourceFactory is a function tha t can be used to create a Datasource from Datasource Properties
type DatasourceFactory func(*Properties) Datasource

// CreateDatasource : this function returns a DatasourceFactory based on Properties type (POSTGRES, MONGO or EDGEDB)
func CreateDatasource(properties *Properties) (DatasourceFactory, error) {
	if properties != nil {
		switch properties.Type {
		case POSTGRES:
			return func(properties *Properties) Datasource {
				return &PgDatasource{
					WithProperties{
						Properties: properties,
					},
				}
			}, nil
		case MONGO:
			return func(properties *Properties) Datasource {
				return &MongoDatasource{
					WithProperties{
						Properties: properties,
					},
				}
			}, nil
		case EDGEDB:
			return func(properties *Properties) Datasource {
				return &EdgeDatasource{
					WithProperties{
						Properties: properties,
					},
				}
			}, nil
		default:
			return nil, errors.New(fmt.Sprintf("TYPE %s NOT SUPPORTED YET", properties.Type))
		}
	}
	return nil, errors.New("DATASOURCE PROPERTIES ARE MISSING")
}
