package middleware

import (
	"github.com/LeonLow97/internal/config"
	"github.com/LeonLow97/internal/pkg/cache"
)

type Middleware struct {
	cfg      config.Config
	appCache cache.Cache
}

func NewMiddleware(cfg config.Config, appCache cache.Cache) *Middleware {
	return &Middleware{
		cfg:      cfg,
		appCache: appCache,
	}
}
