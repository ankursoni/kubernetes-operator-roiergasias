all: build docker-workflow-build

build: build-linux build-osx
build-linux:
	GOOS=linux GOARCH=amd64 go build -o linux/roiergasias .
	chmod +x linux/roiergasias
build-osx:
	GOOS=darwin GOARCH=amd64 go build -o osx/roiergasias .
	chmod +x osx/roiergasias

docker-workflow-build: build-linux
	docker build -t docker.io/ankursoni/roiergasias-operator:workflow .
docker-workflow-push: docker-workflow-build
	docker push docker.io/ankursoni/roiergasias-operator:workflow

docker-examples: build-linux docker-example-local docker-example-aws
docker-example-local:
	docker build -t roiergasias:local -f Dockerfile-local .
	docker tag roiergasias:local docker.io/ankursoni/roiergasias:local
	docker push docker.io/ankursoni/roiergasias:local
docker-example-aws:
	docker build -t roiergasias:aws -f Dockerfile-aws .
	docker tag roiergasias:aws docker.io/ankursoni/roiergasias:aws
	docker push docker.io/ankursoni/roiergasias:aws
