package utils

type CtxKey string

const (
	EntitiesCtxKey    CtxKey = "entities-ctx-key"
	ControllersCtxKey CtxKey = "controllers-ctx-key"
	DBCtxKey          CtxKey = "db-ctx-key"
)
