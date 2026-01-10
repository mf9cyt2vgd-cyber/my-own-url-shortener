package update

import (
	"encoding/json"
	url_validator "my_own_shortener/internal/http-server/url-validator"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type Request struct {
	NewURL string `json:"new_url"`
	Alias  string `json:"alias"`
}
type UrlUpdater interface {
	Update(alias string, newURL string) error
}

func NewUpdateHandler(logger *logrus.Logger, updater UrlUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.handlers.update"
		var req Request
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Errorf("%s:\n\terror decoding json: %s", op, err)
			return
		}
		if !url_validator.ValidateUrl(req.NewURL, 2*time.Second) {
			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).Encode(map[string]string{"result": "invalid url"})
			logger.Errorf("%s:\n\terror validating url", op)
			return
		}
		err = updater.Update(req.Alias, req.NewURL)
		if err != nil {
			logger.Errorf("%s:\n\terror saving url: %s", op, err)
			return
		}
		err = json.NewEncoder(w).Encode(map[string]string{"result": "successful save"})
		if err != nil {
			logger.Errorf("%s:\n\terror encoding json: %s", op, err)
			return
		}
		logger.Infof("%s:\n\tupdate url %s by alias %s successful", op, req.NewURL, req.Alias)
	}
}
