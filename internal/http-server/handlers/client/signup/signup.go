package signup

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"io"
	"log/slog"
	"net/http"
	resp "notify-server/internal/lib/api/response"
	"notify-server/internal/lib/sl"
	"strings"
)

type Request struct {
	Client string `json:"client,omitempty" validate:"required"`
}

type Response struct {
	resp.Response
	Client   string `json:"client,omitempty"`
	Register string `json:"signup_status,omitempty"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=ClientSignup
type ClientSignup interface {
	SaveClient(client string) error
}

func New(log *slog.Logger, clientSaver ClientSignup) http.HandlerFunc {
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

		if client == "" {
			log.Error("Field client is empty on request")

			render.JSON(w, r, resp.Error("Failed to signup client. Field client is empty!"))

			return
		}

		err = clientSaver.SaveClient(client)
		if err != nil {
			if strings.Contains(fmt.Sprint(err), "duplicate key value violates unique constraint \"client_pku\"") {
				render.JSON(w, r, resp.Error("Client already registered"))

				return
			}

			log.Error("failed to signup client", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to signup client"))

			return
		}

		log.Info("signup client success", slog.String("client", req.Client))

		responseOK(w, r, client, "Registered")
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, client string, status string) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Client:   client,
		Register: status,
	})
}
