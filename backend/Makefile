gen:
	@protoc \
		--proto_path=proto "proto/account.proto" \
		--go_out=pkg/api/account \
		--go_opt=paths=source_relative \
		--go-grpc_out=pkg/api/account \
		--go-grpc_opt=paths=source_relative