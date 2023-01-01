package rest

import (
	"net/http"

	"github.com/go-chi/render"
	validation "github.com/go-ozzo/ozzo-validation"
	"go.opentelemetry.io/otel"
)

const otelName = "github.com/MarioCarrion/todo-api/internal/rest"

type ErrorResponse struct {
	Error       string            `json:"error"`
	Validations validation.Errors `json:"validations,omitempty"`
}

func renderErrorResponse(w http.ResponseWriter, r *http.Request, msg string, err error) {
	_, span := otel.Tracer(otelName).Start(r.Context(), "renerErrorResponse")
	defer span.End()
	span.RecordError(err)

	render.Status(r, http.StatusInternalServerError)
	render.JSON(w, r, &ErrorResponse{
		Error: msg,
	})
}

func renderResponse(w http.ResponseWriter, r *http.Request, res interface{}, status int) {
	render.Status(r, status)
	render.JSON(w, r, res)
}
