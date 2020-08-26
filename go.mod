module project

go 1.13

require (
	github.com/gorilla/mux v1.8.0 // indirect
	internal/pkg/handler v1.0.0
)

replace internal/pkg/handler => ./internal/pkg/handler
