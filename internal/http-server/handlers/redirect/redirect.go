package redirect

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type UrlGetter interface {
	Get(alias string) (string, error)
}

func NewRedirectHandler(logger *logrus.Logger, getter UrlGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.handlers.redirect.New"
		alias := chi.URLParam(r, "alias")
		urlToGet, err := getter.Get(alias)
		if err != nil {
			logger.Errorf("%s:\n\terror geting url by alias %s: %w", op, alias, err)
		}
		http.Redirect(w, r, urlToGet, http.StatusFound)
		logger.Infof("%s:\n\tredirect to %s by alias %s complete", op, urlToGet, alias)
	}
}
