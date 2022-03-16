package web

import (
	"context"
	"github.com/hsedjame/webfast/utils"
	"log"
	"net/http"
)

func PostPutMethodHandler(defaultModel Request, modelKey interface{}, errorHandler ErrorHandler) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(wr http.ResponseWriter, rq *http.Request) {

			wr.Header().Add("Content-Type", "application/json")

			if err := utils.FromJson(&defaultModel, rq.Body); err != nil {
				wr.WriteHeader(http.StatusBadRequest)
				_ = errorHandler(err, wr)
				return
			} else if err := defaultModel.Validate(func(request Request) error {
				return utils.IsValid(request)
			}); err != nil {
				wr.WriteHeader(http.StatusBadRequest)
				_ = errorHandler(err, wr)
				return
			}

			ctx := context.WithValue(rq.Context(), modelKey, defaultModel)

			rq = rq.WithContext(ctx)

			next.ServeHTTP(wr, rq)

			return
		})
	}
}

func LoggingMiddleware(logger *log.Logger) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(wr http.ResponseWriter, rq *http.Request) {
			logger.Println("---------------------------------")
			logger.Printf(" ----> [ %s %s ]", rq.Method, rq.URL.Path)
			next.ServeHTTP(wr, rq)
			logger.Println("_________________________________")
		})
	}
}
