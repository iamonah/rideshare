.PHONY: gen-proto

gen-proto:
	@protoc \
	  --proto_path=shared/proto \
	  --go_out=shared/ --go_opt=module=github.com/iamonah/rideshare/shared \
	  --go-grpc_out=shared/ --go-grpc_opt=module=github.com/iamonah/rideshare/shared \
	  shared/proto/*.proto




.PHONY: proto-all
proto-all: SERVICE_NAME=tripservice
proto-all: proto