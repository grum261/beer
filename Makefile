gen-users:
	protoc --go_opt=paths=source_relative --go_out=. \
	--go-grpc_opt=paths=source_relative --go-grpc_out=. \
	--grpc-gateway_out=. --grpc-gateway_opt=paths=source_relative \
	proto/userpb/*.proto

gen-breweries:
	protoc --go_opt=paths=source_relative --go_out=. \
	--go-grpc_opt=paths=source_relative --go-grpc_out=. \
	--grpc-gateway_out=. --grpc-gateway_opt=paths=source_relative \
	proto/brewerypb/*.proto