package std

import (
	"context"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/hsedjame/webfast/data"
	"github.com/hsedjame/webfast/utils"
	"github.com/hsedjame/webfast/web"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"
)

type AppCtx = context.Context

type App struct {
	Logger     *log.Logger
	server     *http.Server
	ctx        AppCtx
	classpath  string
	properties *Properties
	datasource data.Datasource
}

// New : CreateDatasource new app
func New(logger *log.Logger, options ...Options) (*App, error) {

	app := &App{
		Logger: logger,
		ctx:    context.Background(),
	}

	if err := app.setClasspath(); err != nil {
		return nil, err
	}

	return app.
		addToContext(utils.EntitiesCtxKey, []interface{}{}).
		addToContext(utils.ControllersCtxKey, []interface{}{}).
		loadProperties().
		WithOptions(options...), nil
}

//Run launches application
func (app App) Run() {

	// Start web server
	app.configureWebServer()

	// Launch the application
	go func() {
		app.Logger.Fatal(app.server.ListenAndServe())
	}()

	// Shutdown the application
	app.shutDownGracefully()
}

// WithOptions : Add options to application
func (app *App) WithOptions(options ...Options) *App {
	for _, option := range options {
		option(app)
	}
	return app
}

// WithContext : set app context
func (app *App) WithContext(ctx AppCtx) *App {
	app.ctx = ctx
	return app
}

// Context : get app context
func (app App) Context() AppCtx {
	return app.ctx
}

// addToContext adds new value to the application context
func (app *App) addToContext(key utils.CtxKey, value interface{}) *App {
	return app.WithContext(context.WithValue(app.ctx, key, value))
}

// setClasspath set the application classpath
func (app *App) setClasspath() error {
	if rootDir, err := os.Getwd(); err != nil {
		return err
	} else {
		app.classpath = fmt.Sprintf("%s/%s", rootDir, utils.ResourcesLocation)
		app.Logger.Printf("Application classpath configured to : %s \n", app.classpath)
		return nil
	}
}

// configureWebServer configure application web server
func (app *App) configureWebServer() {

	router := mux.NewRouter()

	for _, controller := range app.Context().Value(utils.ControllersCtxKey).([]interface{}) {

		ctrl := (controller).(web.RestController)

		subRouter := router.PathPrefix(ctrl.Path()).Subrouter()

		if ctrl.MiddleWare != nil {
			subRouter.Use(ctrl.MiddleWare)
		}

		subRouter.Use(web.LoggingMiddleware(app.Logger))

		for _, endpoint := range ctrl.Endpoints() {

			var handler http.Handler
			if endpoint.Method() == web.POST || endpoint.Method() == web.PUT {
				model := endpoint.EmptyRequestBody()
				key := endpoint.ModelKey()
				if model != nil && key != nil {
					eh := ctrl.ErrorHandler()
					if eh == nil {
						eh = func(err error, writer http.ResponseWriter) error {
							return utils.ToJson(utils.AppError{Message: err.Error()}, writer)
						}
					}

					handler = web.PostPutMethodHandler(model, key, eh)(endpoint.Handler())
				}
			} else {
				handler = endpoint.Handler()
			}
			subRouter.Handle(endpoint.Path(), handler).Methods(endpoint.Method())
		}
	}

	/* Configure CORS */
	var opts []handlers.CORSOption
	cors := app.properties.Cors
	origins := cors.AllowedOrigins
	headers := cors.AllowedHeaders
	methods := cors.AllowedMethods

	if origins != "" {
		opts = append(opts, handlers.AllowedOrigins(strings.Split(origins, ",")))
	}
	if headers != "" {
		opts = append(opts, handlers.AllowedHeaders(strings.Split(headers, ",")))
	}
	if methods != "" {
		opts = append(opts, handlers.AllowedHeaders(strings.Split(methods, ",")))
	}

	corsHandlers := handlers.CORS(opts...)

	/* Configure server */
	port := app.properties.Server.Port
	if port == 0 {
		port = 8080
	}

	app.server = &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      corsHandlers(router),
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

}

// loadProperties loads the application properties
func (app *App) loadProperties() *App {
	// Load properties
	properties := AppDefaultProperties()
	if err := properties.Load(app.classpath); err != nil {
		app.Logger.Fatal(err)
	}
	app.properties = properties

	// Start database
	var datasource data.Datasource

	dbProperties := properties.Datasource
	if source, err := data.CreateDatasource(dbProperties); err != nil {
		app.Logger.Fatal(err)
	} else {
		datasource = source(dbProperties)
		app.datasource = datasource

		if datasource.CanStart() {
			if ctx, err := datasource.Start(app.Context()); err != nil {
				app.Logger.Fatal(err)
			} else {
				app.WithContext(ctx)
				app.Logger.Println("Database successfully connected")
			}
		}
	}
	return app
}

func (app App) shutDownGracefully() {

	// Create a channel to listen OS signals
	osSignalsChannel := make(chan os.Signal)

	// Send a message to the channel when
	//  - interruption occurs
	//  - os is killed
	signal.Notify(osSignalsChannel, os.Interrupt)
	signal.Notify(osSignalsChannel, os.Kill)

	// Wait for new signal
	_ = <-osSignalsChannel

	app.Logger.Println(" ### Arrêt du serveur ....")

	deadline, cancel := context.WithTimeout(app.Context(), 30*time.Second)

	defer cancel()

	// stop datasource
	if app.datasource != nil {
		app.datasource.Stop(app.Context())
	}

	// Shutdown the server
	app.Logger.Fatal(app.server.Shutdown(deadline))

}
