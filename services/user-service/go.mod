module github.com/franzego/user-service

replace github.com/franzego/redisclient => ../pkg/redisclient

go 1.25.0

require (
	github.com/franzego/redisclient v0.0.0-00010101000000-000000000000
	github.com/go-playground/validator/v10 v10.27.0
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/gorilla/mux v1.8.1
	github.com/jackc/pgx/v5 v5.7.6
	github.com/joho/godotenv v1.5.1
	github.com/redis/go-redis/v9 v9.14.0
	golang.org/x/crypto v0.41.0
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/gabriel-vasile/mimetype v1.4.8 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	golang.org/x/net v0.43.0 // indirect
	golang.org/x/sync v0.16.0 // indirect
	golang.org/x/sys v0.35.0 // indirect
	golang.org/x/text v0.28.0 // indirect
)
