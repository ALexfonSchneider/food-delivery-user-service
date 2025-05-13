
grpc-gen:
	sudo rm -rf gen/grpc/go && sudo mkdir -p gen/grpc/go && protoc -I food-delivery-proto --go_out=gen/grpc/go --go-grpc_out=gen/grpc/go --go_opt=module=github.com/ALexfonSchneider/food-delivery-proto --go-grpc_opt=module=github.com/ALexfonSchneider/food-delivery-proto ./food-delivery-proto/user/*.proto