package rest

import (
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi"
	"gopkg.in/yaml.v2"
)

//go:generate go run ../../cmd/openapi_gen/main.go -path .
//go:generate oapi-codegen -package openapi3 -old-config-style -generate types  -o ../../pkg/openapi3/task_types.gen.go openapi3.yaml
//go:generate oapi-codegen -package openapi3 -old-config-style -generate client -o ../../pkg/openapi3/client.gen.go     openapi3.yaml

func NewOpenAPI3() openapi3.T {
	swagger := openapi3.T{
		OpenAPI: "3.0.0",
		Info: &openapi3.Info{
			Title:       "Ntube",
			Description: "Ntube REST APIs",
			Version:     "0.0.1",
			License: &openapi3.License{
				Name: "MIT",
				URL:  "https://opensource.org/licenses/MIT",
			},
			Contact: &openapi3.Contact{
				URL: "https://github.com/nei7/ntube",
			},
		},
		Servers: openapi3.Servers{
			&openapi3.Server{
				Description: "Rest server",
				URL:         "http://127.0.0.1:3001",
			},
			&openapi3.Server{
				Description: "Video server",
				URL:         "http://127.0.0.1:3002",
			},
		},
		Components: &openapi3.Components{},
	}

	swagger.Components.Schemas = openapi3.Schemas{
		"User": openapi3.NewSchemaRef("",
			openapi3.NewObjectSchema().
				WithProperty("id", openapi3.NewUUIDSchema()).
				WithProperty("description", openapi3.NewStringSchema().WithNullable()).
				WithProperty("avatar", openapi3.NewStringSchema().WithNullable()).
				WithProperty("followers", openapi3.NewIntegerSchema()).
				WithProperty("created_at", openapi3.NewDateTimeSchema()).
				WithProperty("email", openapi3.NewStringSchema()),
		),
		"Video": openapi3.NewSchemaRef("", openapi3.NewObjectSchema().
			WithProperty("id", openapi3.NewUUIDSchema()).
			WithProperty("description", openapi3.NewStringSchema()).
			WithProperty("title", openapi3.NewStringSchema()).
			WithProperty("thumbnail", openapi3.NewStringSchema()).
			WithProperty("path", openapi3.NewStringSchema()).
			WithProperty("uploaded_at", openapi3.NewStringSchema().
				WithFormat("date-time").
				WithNullable(),
			).WithPropertyRef("owner", &openapi3.SchemaRef{
			Ref: "#/components/schemas/User",
		})),
	}

	swagger.Components.RequestBodies = openapi3.RequestBodies{
		"AuthRequest": &openapi3.RequestBodyRef{
			Value: openapi3.NewRequestBody().
				WithDescription("Request used for creating an account.").
				WithRequired(true).
				WithJSONSchema(openapi3.NewSchema().
					WithProperty("email", openapi3.NewStringSchema().
						WithMinLength(1)).
					WithProperty("password", openapi3.NewStringSchema().
						WithMinLength(1))),
		},
		"SearchVideoRequest": &openapi3.RequestBodyRef{
			Value: openapi3.NewRequestBody().
				WithDescription("Request used for searching a video.").
				WithRequired(true).
				WithJSONSchema(openapi3.NewSchema().
					WithProperty("title", openapi3.NewStringSchema().
						WithMinLength(1)).
					WithProperty("from", openapi3.NewInt64Schema().
						WithDefault(0)).
					WithProperty("size", openapi3.NewInt64Schema().
						WithDefault(10))),
		},
		"RenewTokenRequest": &openapi3.RequestBodyRef{
			Value: openapi3.NewRequestBody().
				WithDescription("Request used for renewing jwt token.").
				WithRequired(true).
				WithJSONSchema(openapi3.NewSchema().
					WithProperty("token", openapi3.NewStringSchema().
						WithMinLength(1))),
		},
	}

	swagger.Components.Responses = openapi3.Responses{
		"ErrorResponse": &openapi3.ResponseRef{
			Value: openapi3.NewResponse().
				WithDescription("Response when errors happen.").
				WithContent(openapi3.NewContentWithJSONSchema(openapi3.NewSchema().
					WithProperty("error", openapi3.NewStringSchema()))),
		},
		"SearchVideoResponse": &openapi3.ResponseRef{
			Value: openapi3.NewResponse().
				WithDescription("Response returned back after searching for any videos.").
				WithContent(openapi3.NewContentWithJSONSchema(openapi3.NewSchema().
					WithPropertyRef("tasks", &openapi3.SchemaRef{
						Value: &openapi3.Schema{
							Type: "array",
							Items: &openapi3.SchemaRef{
								Ref: "#/components/schemas/Video",
							},
						},
					}).
					WithProperty("total", openapi3.NewInt64Schema()))),
		},
		"UserResponse": &openapi3.ResponseRef{
			Value: openapi3.NewResponse().
				WithDescription("Login response.").
				WithContent(openapi3.NewContentWithJSONSchema(openapi3.NewSchema().
					WithPropertyRef("user", &openapi3.SchemaRef{
						Ref: "#/components/schemas/User",
					}).
					WithProperty("access_token", openapi3.NewStringSchema()).
					WithProperty("refresh_token", openapi3.NewStringSchema()),
				)),
		},
	}

	swagger.Paths = openapi3.Paths{
		"/signup": &openapi3.PathItem{
			Post: &openapi3.Operation{
				OperationID: "Signup",
				RequestBody: &openapi3.RequestBodyRef{
					Ref: "#/components/requestBodies/AuthRequest",
				},
				Responses: openapi3.Responses{
					"400": &openapi3.ResponseRef{
						Ref: "#/components/responses/ErrorResponse",
					},
					"500": &openapi3.ResponseRef{
						Ref: "#/components/responses/ErrorResponse",
					},
					"201": &openapi3.ResponseRef{
						Ref: "#/components/responses/UserResponse",
					},
				},
			},
		},
		"/login": &openapi3.PathItem{
			Post: &openapi3.Operation{
				OperationID: "Login",
				RequestBody: &openapi3.RequestBodyRef{
					Ref: "#/components/requestBodies/AuthRequest",
				},
				Responses: openapi3.Responses{
					"400": &openapi3.ResponseRef{
						Ref: "#/components/responses/ErrorResponse",
					},
					"500": &openapi3.ResponseRef{
						Ref: "#/components/responses/ErrorResponse",
					},
					"200": &openapi3.ResponseRef{
						Ref: "#/components/responses/UserResponse",
					},
				},
			},
		},
		"/videos/search": &openapi3.PathItem{
			Post: &openapi3.Operation{
				OperationID: "SearchVideos",
				RequestBody: &openapi3.RequestBodyRef{
					Ref: "#/components/requestBodies/SearchVideoRequest",
				},
				Responses: openapi3.Responses{
					"200": &openapi3.ResponseRef{
						Ref: "#/components/responses/SearchVideoResponse",
					},
					"400": &openapi3.ResponseRef{
						Ref: "#/components/responses/ErrorResponse",
					},
					"500": &openapi3.ResponseRef{
						Ref: "#/components/responses/ErrorResponse",
					},
				},
			},
		},
	}

	return swagger
}

func RegisterOpenAPI(router *chi.Mux) {
	swagger := NewOpenAPI3()

	router.Get("/openapi3.json", func(w http.ResponseWriter, r *http.Request) {
		renderResponse(w, r, &swagger, http.StatusOK)
	})

	router.Get("/openapi3.yaml", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/x-yaml")

		data, _ := yaml.Marshal(&swagger)

		_, _ = w.Write(data)

		w.WriteHeader(http.StatusOK)
	})
}
