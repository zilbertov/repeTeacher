package docs

import "github.com/swaggo/swag"

const docTemplate = `{
  "swagger": "2.0",
  "info": {
    "title": "RepeTeacher Users API",
    "description": "Students, profile, chats and notifications service.",
    "version": "1.0"
  },
  "host": "localhost:8081",
  "basePath": "/",
  "schemes": ["http"],
  "paths": {
    "/health": {"get": {"summary": "Healthcheck", "responses": {"200": {"description": "OK"}}}},
    "/api/auth/login": {"post": {"summary": "Demo JWT login", "responses": {"200": {"description": "OK"}}}},
    "/api/students": {
      "get": {"summary": "List students", "responses": {"200": {"description": "OK"}}},
      "post": {"summary": "Create student", "responses": {"201": {"description": "Created"}}}
    },
    "/api/students/{id}": {
      "get": {"summary": "Get student", "responses": {"200": {"description": "OK"}}},
      "put": {"summary": "Update student", "responses": {"200": {"description": "OK"}}},
      "delete": {"summary": "Delete student", "responses": {"204": {"description": "Deleted"}}}
    },
    "/api/students/{id}/accept": {"post": {"summary": "Accept student request", "responses": {"200": {"description": "OK"}}}},
    "/api/students/{id}/archive": {"post": {"summary": "Archive student", "responses": {"200": {"description": "OK"}}}},
    "/api/students/{id}/notes": {"post": {"summary": "Update student notes", "responses": {"200": {"description": "OK"}}}},
    "/api/profile": {
      "get": {"summary": "Get tutor profile", "responses": {"200": {"description": "OK"}}},
      "put": {"summary": "Update tutor profile", "responses": {"200": {"description": "OK"}}}
    },
    "/api/profile/password": {"post": {"summary": "Change password", "responses": {"200": {"description": "OK"}}}},
    "/api/settings/notifications": {
      "get": {"summary": "Get notification settings", "responses": {"200": {"description": "OK"}}},
      "put": {"summary": "Update notification settings", "responses": {"200": {"description": "OK"}}}
    },
    "/api/chats": {
      "get": {"summary": "List chats", "responses": {"200": {"description": "OK"}}},
      "post": {"summary": "Create chat", "responses": {"201": {"description": "Created"}}}
    },
    "/api/chats/{id}/messages": {
      "get": {"summary": "List messages", "responses": {"200": {"description": "OK"}}},
      "post": {"summary": "Send message", "responses": {"201": {"description": "Created"}}}
    },
    "/api/notifications": {"get": {"summary": "List notifications", "responses": {"200": {"description": "OK"}}}},
    "/api/tutors": {
      "get": {"summary": "List demo tutors", "responses": {"200": {"description": "OK"}}},
      "post": {"summary": "Create demo tutor", "responses": {"201": {"description": "Created"}}}
    },
    "/api/tutors/{id}": {"get": {"summary": "Get demo tutor", "responses": {"200": {"description": "OK"}}}},
    "/api/tutors/{id}/notes": {"post": {"summary": "Update tutor notes", "responses": {"200": {"description": "OK"}}}},
    "/api/notifications/{id}/read": {"post": {"summary": "Mark notification read", "responses": {"200": {"description": "OK"}}}},
    "/api/notifications/{id}/approve": {"post": {"summary": "Approve notification", "responses": {"200": {"description": "OK"}}}},
    "/api/notifications/{id}/reject": {"post": {"summary": "Reject notification", "responses": {"200": {"description": "OK"}}}}
  }
}`

var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8081",
	BasePath:         "/",
	Schemes:          []string{"http"},
	Title:            "RepeTeacher Users API",
	Description:      "Students, profile, chats and notifications service.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
