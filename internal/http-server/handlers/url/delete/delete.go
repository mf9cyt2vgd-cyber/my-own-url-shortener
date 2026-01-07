package save

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

type Request struct {
	Alias string `json:"alias"`
}

type UrlDeleter interface {
	Delete(alias string) error
}

func NewDeleteHandler(logger *logrus.Logger, deleter UrlDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.handlers.save"
		var req Request
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Errorf("error decoding json: %s", err)
			return
		}
		err = deleter.Delete(req.Alias)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Errorf("error deleting alias %s", req.Alias)
			return
		}
		err = json.NewEncoder(w).Encode(map[string]string{"result": "successful delete"})
		logger.Info("delete successful")
	}
}
