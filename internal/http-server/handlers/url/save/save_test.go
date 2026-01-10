package save

import (
	"bytes"
	"encoding/json"
	"io"
	"my_own_shortener/internal/http-server/handlers/url/save/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSaveHandler_Success(t *testing.T) {
	logger := logrus.New()
	logger.SetOutput(io.Discard)

	mockSaver := mocks.NewUrlSaver(t)
	handler := NewSaveHandler(logger, mockSaver)

	mockSaver.On("Save", mock.Anything, "https://example.com", "custom-alias").
		Return(nil).
		Once()

	requestBody := map[string]string{
		"url":   "https://example.com",
		"alias": "custom-alias",
	}
	body, _ := json.Marshal(requestBody)

	req := httptest.NewRequest(http.MethodPost, "/save", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "successful save", response["result"])
	mockSaver.AssertExpectations(t)
}
