package save

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"my_own_shortener/internal/random"

	"my_own_shortener/internal/http-server/url-validator"

	"github.com/sirupsen/logrus"
)

type Request struct {
	URL   string `json:"url"`
	Alias string `json:"alias"`
}

const aliasLength = 6

//go:generate mockery --name=UrlSaver --output=mocks --outpkg=mocks --case=underscore
type UrlSaver interface {
	Save(ctx context.Context, urlToSave, alias string) error
}

func NewSaveHandler(logger *logrus.Logger, saver UrlSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.handlers.save"
		var req Request
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Errorf("%s:\n\terror decoding json: %s", op, err)
			return
		}
		if req.Alias == "" {
			req.Alias = random.NewRandomString(aliasLength)
		}
		if !url_validator.ValidateUrl(req.URL, 2*time.Second) {
			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).Encode(map[string]string{"result": "invalid url"})
			logger.Errorf("%s:\n\terror validating url", op)
			return
		}
		err = saver.Save(r.Context(), req.URL, req.Alias)
		if err != nil {
			logger.Errorf("%s:\n\terror saving url: %s", op, err)
			err = json.NewEncoder(w).Encode(map[string]string{"result": "error saving url"})
			return
		}
		err = json.NewEncoder(w).Encode(map[string]string{"result": "successful save"})
		if err != nil {
			logger.Errorf("%s:\n\terror encoding json: %s", op, err)
			return
		}
		logger.Infof("%s:\n\tsave url %s by alias %s successful", op, req.URL, req.Alias)
	}
}
