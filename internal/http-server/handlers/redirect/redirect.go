package redirect

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type UrlGetter interface {
	Get(ctx context.Context, alias string) (string, error)
}

func NewRedirectHandler(logger *logrus.Logger, getter UrlGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.handlers.redirect.New"
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		alias := chi.URLParam(r, "alias")
		urlToGet, err := getter.Get(ctx, alias)
		if err != nil {
			logger.Errorf("%s:\n\terror geting url by alias %s: %w", op, alias, err)
			return
		}
		http.Redirect(w, r, urlToGet, http.StatusFound)
		logger.Infof("%s:\n\tredirect to %s by alias %s complete", op, urlToGet, alias)
	}
}
