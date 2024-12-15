package check

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
	Client   string `json:"client,omitempty"`
	IsActive bool   `json:"isActive,omitempty"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=ClientCheckSignup
type ClientCheckSignup interface {
	ClientRegistered(client string) (bool, error)
}

func New(log *slog.Logger, clientChecker ClientCheckSignup) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.clientCheckSignup.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		client := chi.URLParam(r, "client")

		if client == "" {
			log.Info("client is empty")

			render.JSON(w, r, resp.Error("client is empty"))

			return
		}

		log.Info("check client on registered", slog.Any("client", client))

		status, err := clientChecker.ClientRegistered(client)

		if err != nil {
			log.Error("failed to check client", sl.Err(err))

			if strings.Contains(fmt.Sprint(err), "no rows in result set") {
				render.JSON(w, r, resp.Error("Client not registered"))

				return
			}

			render.JSON(w, r, resp.Error("Failed check user"))

			return
		}

		responseOK(w, r, client, status)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, client string, status bool) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Client:   client,
		IsActive: status,
	})
}
