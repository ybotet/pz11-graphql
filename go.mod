module github.com/ybotet/pz8-pipelineCICD-go

go 1.25.1

require (
	github.com/99designs/gqlgen v0.17.45
	github.com/golang-jwt/jwt/v5 v5.2.0
	github.com/google/uuid v1.6.0
	github.com/gorilla/mux v1.8.1
	github.com/lib/pq v1.11.2
	github.com/redis/go-redis/v9 v9.18.0
	github.com/sirupsen/logrus v1.9.4
	github.com/vektah/gqlparser/v2 v2.5.11
	github.com/ybotet/pz8-pipelineCICD-go/gen v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.79.1
)

require (
	github.com/agnivade/levenshtein v1.1.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/hashicorp/golang-lru/v2 v2.0.7 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/sosodev/duration v1.2.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	golang.org/x/net v0.52.0 // indirect
	golang.org/x/sys v0.42.0 // indirect
	golang.org/x/text v0.35.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251202230838-ff82c1b0f217 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)

// Si tienes código generado en /gen, mantenemos el replace
replace github.com/ybotet/pz8-pipelineCICD-go/gen => ./gen
