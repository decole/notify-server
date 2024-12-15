package save

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"io"
	"log/slog"
	"net/http"
	"notify-server/internal/lib/sl"
	"notify-server/internal/storage/postgres"

	resp "notify-server/internal/lib/api/response"
)

type Request struct {
	Client  string `json:"client,omitempty"`
	Message string `json:"message" validate:"required"`
}

type Response struct {
	resp.Response
	Client  string `json:"client,omitempty"`
	Message string `json:"message,omitempty"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=SaveNotify
type NotifySaver interface {
	SaveNotify(client string, message string) error
	GetActiveUsers() ([]postgres.Client, error)
}

func New(log *slog.Logger, notifySaver NotifySaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.save.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			// Такую ошибку встретим, если получили запрос с пустым телом.
			// Обработаем её отдельно
			log.Error("request body is empty")

			render.JSON(w, r, resp.Error("empty request"))

			return
		}
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request"))

			// прерываем запрос
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}

		client := req.Client

		if req.Message == "" {
			log.Error("Field message is empty on request")

			render.JSON(w, r, resp.Error("Failed to save notify. Field message is empty!"))

			return
		}

		if client == "" {
			users, err := notifySaver.GetActiveUsers()

			for _, u := range users {
				singleClient := fmt.Sprint(u)
				err = notifySaver.SaveNotify(singleClient, req.Message)

				if err != nil {
					log.Error("failed to save notify", sl.Err(err))

					render.JSON(w, r, resp.Error("failed to save notify"))

					return
				}
			}

			responseOK(w, r, "all", req.Message)

			return
		}

		err = notifySaver.SaveNotify(client, req.Message)
		if err != nil {
			log.Error("failed to save notify", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to save notify"))

			return
		}

		log.Info("notify added", slog.String("message", req.Message), slog.String("client", req.Client))

		responseOK(w, r, client, req.Message)
	}
}

func saveSingle(client string, message string) {

}

func responseOK(w http.ResponseWriter, r *http.Request, client string, message string) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Client:   client,
		Message:  message,
	})
}
