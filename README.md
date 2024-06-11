```
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
go get github.com/grpc-ecosystem/grpc-gateway/v2/runtime
go get github.com/grpc-ecosystem/grpc-gateway/v2/utilities
export PATH="$PATH:$(go env GOPATH)/bin"
go get github.com/jackc/pgx/v4
go get google.golang.org/grpc
go get google.golang.org/grpc/reflection

protoc -I . --go_out=. --go-grpc_out=. --grpc-gateway_out=. --grpc-gateway_opt=logtostderr=true email.proto

invoke request
{
    "recipient": "recipient1",
    "subject": "subject1",
    "body": "body1"
}

POST localhost:8081/email
raw JSON
{
    "recipient": "recipient2",
    "subject": "subject2",
    "body": "body2"
}
```