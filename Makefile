proto:
	protoc --go_out=./query --go_opt=paths=source_relative \
	 --go-grpc_out=./query --go-grpc_opt=paths=source_relative \
	  query.proto
envoy:
	docker run --rm -it -v $(CURDIR)/envoy.yaml:/envoy.yaml -p 80:80 envoyproxy/envoy:dev -c envoy.yaml