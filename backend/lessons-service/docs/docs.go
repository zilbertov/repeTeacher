package docs

import "github.com/swaggo/swag"

const docTemplate = `{
  "swagger": "2.0",
  "info": {
    "title": "RepeTeacher Lessons API",
    "description": "Calendar lessons and lesson files service.",
    "version": "1.0"
  },
  "host": "localhost:8082",
  "basePath": "/",
  "schemes": ["http"],
  "paths": {
    "/health": {"get": {"summary": "Healthcheck", "responses": {"200": {"description": "OK"}}}},
    "/api/lessons": {
      "get": {"summary": "List lessons", "responses": {"200": {"description": "OK"}}},
      "post": {"summary": "Create lesson", "responses": {"201": {"description": "Created"}}}
    },
    "/api/lessons/{id}": {
      "get": {"summary": "Get lesson", "responses": {"200": {"description": "OK"}}},
      "put": {"summary": "Update lesson", "responses": {"200": {"description": "OK"}}}
    },
    "/api/lessons/{id}/reschedule": {"post": {"summary": "Reschedule lesson", "responses": {"200": {"description": "OK"}}}},
    "/api/lessons/{id}/cancel": {"post": {"summary": "Cancel lesson", "responses": {"200": {"description": "OK"}}}},
    "/api/lessons/{id}/files": {"post": {"summary": "Add lesson file", "responses": {"201": {"description": "Created"}}}}
  }
}`

var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8082",
	BasePath:         "/",
	Schemes:          []string{"http"},
	Title:            "RepeTeacher Lessons API",
	Description:      "Calendar lessons and lesson files service.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
