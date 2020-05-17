package tools

//go:generate protoc -I. -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway --go_out=plugins=grpc:. --grpc-gateway_out=logtostderr=true:. --swagger_out=logtostderr=true:../swagger user/user.proto

// import (
// 	_ "github.com/golang/protobuf/protoc-gen-go"
// 	_ "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway"
// 	_ "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger"
// )
