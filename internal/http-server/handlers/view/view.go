package view

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	resp "notify-server/internal/lib/api/response"
	"notify-server/internal/lib/sl"
	"strings"
)

type Response struct {
	resp.Response
	Client  string `json:"client,omitempty"`
	Message string `json:"message,omitempty"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=NotifyGetter
type NotifyGetter interface {
	GetNotify(client string) (string, error)
}

func New(log *slog.Logger, notifyGetter NotifyGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.view.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		client := chi.URLParam(r, "client")

		if client == "" {
			log.Info("client is empty")

			render.JSON(w, r, resp.Error("invalid request"))

			return
		}

		log.Info("get notify by client", slog.Any("client", client))

		message, err := notifyGetter.GetNotify(client)

		if err != nil {
			log.Error("failed to get notify", sl.Err(err))

			if strings.Contains(fmt.Sprint(err), "no rows in result set") {
				render.JSON(w, r, resp.Error("Notifications is empty"))

				return
			}

			render.JSON(w, r, resp.Error("Failed to get notify"))

			return
		}

		responseOK(w, r, client, message)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, client string, message string) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Client:   client,
		Message:  message,
	})
}
