module github.com/zilbertov/repe-teacher-users-service

go 1.26

require (
	github.com/DATA-DOG/go-sqlmock v1.5.2
	github.com/lib/pq v1.10.9
	github.com/swaggo/http-swagger v1.3.4
	github.com/swaggo/swag v1.16.6
	github.com/zilbertov/repe-teacher-common v0.0.0
	go.uber.org/zap v1.27.0
)

require (
	github.com/KyleBanks/depth v1.2.1
	github.com/go-openapi/jsonpointer v0.19.5
	github.com/go-openapi/jsonreference v0.20.0
	github.com/go-openapi/spec v0.20.6
	github.com/go-openapi/swag v0.19.15
	github.com/josharian/intern v1.0.0
	github.com/mailru/easyjson v0.7.6
	github.com/swaggo/files v0.0.0-20220610200504-28940afbdbfe
	go.uber.org/multierr v1.10.0
	golang.org/x/mod v0.21.0
	golang.org/x/net v0.38.0
	golang.org/x/sync v0.12.0
	golang.org/x/tools v0.24.0
	gopkg.in/yaml.v2 v2.4.0
)

replace github.com/zilbertov/repe-teacher-common => ../common
