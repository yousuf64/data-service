proto:
	protoc --go_out=./query --go_opt=paths=source_relative \
	 --go-grpc_out=./query --go-grpc_opt=paths=source_relative \
	  query.proto