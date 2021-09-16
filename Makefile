SAMPLE_WATCH_DIRS=/churro
GRPC_CERTS_DIR=certs/grpc
DB_CERTS_DIR=certs/db
BUILDDIR=./build
PIPELINE=test
CHURRO_NS=churro
.DEFAULT_GOAL := all

deploy-extension:
	kubectl -n $(PIPELINE) create -f deploy/extension/spacex-geo-extension-deployment.yaml
	kubectl -n $(PIPELINE) create -f deploy/extension/spacex-geo-extension-service.yaml
push:
	docker push docker.io/churrodata/spacex-geo-extension

compile-extension:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=require_unimplemented_servers=false:. --go-grpc_opt=paths=source_relative rpc/extension/spacex-geo-extension.proto
	go build -o build/spacex-geo-extension cmd/spacex-geo-extension/spacex-geo-extension.go
build-extension: compile-extension
	docker build -f ./images/Dockerfile.spacex-geo-extension -t docker.io/churrodata/spacex-geo-extension .

all: build-extension

.PHONY: clean

clean:
	rm $(BUILDDIR)/churro*
	rm /tmp/churro*.log