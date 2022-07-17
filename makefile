genproto:
	protoc --go-grpc_out=plugins=grpc:.  "protos/*.proto"
genprotogrpc:
	protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    protos/*.proto
runserver:
	go run server/server.go
runclient:
	go run client/client.go
.phony: genproto