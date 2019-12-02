help:
	@echo "Please use 'make <target>' where <target> is one of the following:"
	@echo "  serve           to serve the app."
	@echo "  run             to run the app without building."
	@echo "  build           to build the app."
	@echo "  lint            to perform linting."
	@echo "  fmt             to perform formatting."
	@echo "  ci              to run the tests on ci pipeline."
	@echo "  ci-cleanup      to kill & remove all ci containers."
	@echo "  test            to run the tests."
	@echo "  test-init       to run all the tests with integration."

serve:
	docker-compose up -d pigeon

stop:
	docker-compose down

lint:
	golint -set_exit_status=1 `go list ./...`

test:
	go test -mod=vendor `go list ./... | grep -v 'docs'` -race

test-int:
	go test -mod=vendor `go list ./... | grep -v 'docs'` -race -tags=integration

ci:
	docker-compose -f infra/deploy/local/docker-compose.yaml down
	docker-compose -f infra/deploy/local/docker-compose.yaml build pigeon_ci
	docker-compose -f infra/deploy/local/docker-compose.yaml run pigeon_ci ./script/ci.sh
	docker-compose -f infra/deploy/local/docker-compose.yaml down

ci-cleanup:
	docker-compose -f infra/deploy/local/docker-compose.yaml down

run:
	go run -mod=vendor ./cmd/pigeon/main.go

build:
	CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -installsuffix cgo -o ./pigeon ./cmd/pigeon/main.go

fmt:
	go fmt ./...