package middleware

import "github.com/LeonLow97/internal/config"

type Middleware struct {
	cfg config.Config
}

func NewMiddleware(cfg config.Config) *Middleware {
	return &Middleware{
		cfg: cfg,
	}
}
