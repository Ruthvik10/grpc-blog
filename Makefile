build-client:
	@go build -o ./bin ./client

build-server: build-client
	@go build -o ./bin ./server

gen: build-server
	@protoc -Iproto/ --go_opt=module=github.com/Ruthvik10/grpc-blog --go_out=. --go-grpc_opt=module=github.com/Ruthvik10/grpc-blog --go-grpc_out=. proto/*.proto

	