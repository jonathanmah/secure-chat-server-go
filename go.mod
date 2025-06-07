module chatapp

go 1.24.3

require (
	github.com/go-chi/chi/v5 v5.2.1
	github.com/golang-jwt/jwt/v5 v5.2.2
	github.com/gorilla/websocket v1.5.3
	github.com/lib/pq v1.10.9
	golang.org/x/crypto v0.38.0
	golang.org/x/oauth2 v0.30.0
)

require cloud.google.com/go/compute/metadata v0.3.0 // indirect
