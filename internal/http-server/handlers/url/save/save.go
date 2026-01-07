package save

import (
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

type UrlSaver interface {
	Save(urlToSave, alias string) error
}

func NewSaveHandler(logger *logrus.Logger, saver UrlSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.handlers.save"
		var req Request
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Errorf("error decoding json: %s", err)
			return
		}
		if req.Alias == "" {
			req.Alias = random.NewRandomString(aliasLength)
		}
		if !url_validator.ValidateUrl(req.URL, 2*time.Second) {
			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).Encode(map[string]string{"result": "invalid url"})
			logger.Errorf("error validating url")
			return
		}
		err = saver.Save(req.URL, req.Alias)
		if err != nil {
			logger.Errorf("error saving url: %s", err)
			return
		}
		err = json.NewEncoder(w).Encode(map[string]string{"result": "successful save"})
		if err != nil {
			logger.Errorf("error encoding json: %s", err)
			return
		}
		logger.Info("save successful")
	}
}
