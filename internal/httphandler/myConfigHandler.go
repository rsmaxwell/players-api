package httphandler

import (
	"context"
	"net/http"
	"time"

	"github.com/rsmaxwell/players-api/internal/config"
)

type MyConfig struct {
	handler http.Handler
	config  *config.Config
}

func AddConfigContext(handlerToWrap http.Handler, config *config.Config) *MyConfig {
	return &MyConfig{handler: handlerToWrap, config: config}
}

func (h *MyConfig) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	ctx1, cancel := context.WithTimeout(r.Context(), time.Duration(60*time.Second))
	defer cancel()

	ctx2 := context.WithValue(ctx1, ContextConfigKey, h.config)
	r3 := r.WithContext(ctx2)

	h.handler.ServeHTTP(w, r3)
}
