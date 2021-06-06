genG:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative notify/protos/*.proto
clean:
	rm notify/protos/*.go
runs:
	go run notify/server/server.go
runc:
	go run notify/client/client.go