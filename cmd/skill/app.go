package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/VladimirAzanza/alisa_skill/internal/logger"
	"github.com/VladimirAzanza/alisa_skill/internal/models"
	"github.com/VladimirAzanza/alisa_skill/internal/store"
	"go.uber.org/zap"
)

type app struct {
	store store.Store
}

func newApp(s store.Store) *app {
	return &app{store: s}
}

// to run this service:
// curl -X POST http://localhost:8080 -H "Content-Type: application/json" -d '{}'
func (s *app) webhook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodPost {
		logger.Log.Debug("got request with bad method", zap.String("method", r.Method))
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	logger.Log.Debug("decoding request")
	var req models.Request
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		logger.Log.Debug("cannot decode request JSON body", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if req.Request.Type != models.TypeSimpleUtterance {
		logger.Log.Debug("unsupported request type", zap.String("type", req.Request.Type))
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	messages, err := s.store.ListMessages(ctx, req.Session.User.UserID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	text := "Для вас нет новых сообщений."
	if len(messages) > 0 {
		text = fmt.Sprintf("Для вас %d новых сообщений.", len(messages))
	}

	if req.Session.New {
		tz, err := time.LoadLocation(req.Timezone)
		if err != nil {
			logger.Log.Debug("cannot parse timezone")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		now := time.Now().In(tz)
		hour, minute, _ := now.Clock()

		text = fmt.Sprintf("Точное время %d часов, %d минут. %s", hour, minute, text)
	}

	resp := models.Response{
		Response: models.ResponsePayload{
			Text: text,
		},
		Version: "1.0",
	}

	w.Header().Set("Content-Type", "application/json")

	enc := json.NewEncoder(w)
	if err := enc.Encode(resp); err != nil {
		logger.Log.Debug("error encoding response", zap.Error(err))
		return
	}
	logger.Log.Debug("sending HTTP 200 response")
}
