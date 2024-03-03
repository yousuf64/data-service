envoy:
	docker run --rm -it -v $(CURDIR)/envoy.yaml:/envoy.yaml -p 80:80 envoyproxy/envoy:dev -c envoy.yaml
gen-stubs:
	$(foreach file, $(wildcard proto/*.proto), \
			dir=`echo $(basename $(file))`; \
			filename=`echo $$dir | sed 's:.*/::'`; \
			mkdir -p $$dir; \
			protoc --go_out=. --go_opt=paths=source_relative \
				--go-grpc_out=. --go-grpc_opt=paths=source_relative $(file); \
			mkdir -p proto/gen/$$filename; \
			mv proto/$${filename}.pb.go proto/gen/$$filename/$${filename}.pb.go; \
			mv proto/$${filename}_grpc.pb.go proto/gen/$$filename/$${filename}_grpc.pb.go; \
			rmdir proto/$${filename};)
run:
	make envoy; \
	go run ./go-messages-api; \
	go run ./message-query-service --env=message-query-service/env.yaml; \
	go run ./message-query-service --env=message-query-service/env.yaml;
db:
	docker run --name scylla --hostname scylla -d scylladb/scylla;