package rest

import (
	"net/http"

	"github.com/go-chi/render"
	validation "github.com/go-ozzo/ozzo-validation"
	"go.opentelemetry.io/otel"
)

const otelName = "github.com/nei7/ntube/rest"

type ErrorResponse struct {
	Error       string            `json:"error"`
	Validations validation.Errors `json:"validations,omitempty"`
}

func renderErrorResponse(w http.ResponseWriter, r *http.Request, err error, status int) {
	_, span := otel.Tracer(otelName).Start(r.Context(), "renderErrorResponse")
	defer span.End()
	span.RecordError(err)

	render.Status(r, status)
	render.JSON(w, r, &ErrorResponse{
		Error: err.Error(),
	})
}

func renderResponse(w http.ResponseWriter, r *http.Request, res interface{}, status int) {
	render.Status(r, status)
	render.JSON(w, r, res)
}
