package std

import "github.com/hsedjame/webfast/utils"

// Options : Application options
type Options func(app *App)

// Rest : Declare application restControllers
func Rest(controllers []interface{}) Options {
	return func(app *App) {
		restControllers := app.Context().Value(utils.ControllersCtxKey).([]interface{})
		if restControllers == nil {
			restControllers = []interface{}{}
		}
		restControllers = append(restControllers, controllers...)
		app.addToContext(utils.ControllersCtxKey, restControllers)
		app.Logger.Printf("%d RestControllers registered", len(restControllers))
	}
}

// Entities : Declare application entities
func Entities(datas []interface{}) Options {
	return func(app *App) {
		entities := app.Context().Value(utils.EntitiesCtxKey).([]interface{})
		if entities == nil {
			entities = []interface{}{}
		}
		entities = append(entities, datas...)

		app.addToContext(utils.EntitiesCtxKey, entities)

		if err := app.datasource.CreateTables(app.Context()); err != nil {
			app.Logger.Fatal(err)
		}

		app.Logger.Printf("%d Entities registered", len(entities))
	}
}
