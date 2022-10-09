go mod init edgemgmt.com/go-edgemgmt-grpc
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative edgemgmt/edgemgmt.proto

go mod tidy
go run edgemgmt/server/edgemgmt_server.go 
 go run edgemgmt/client/edgemgmt_client.go 