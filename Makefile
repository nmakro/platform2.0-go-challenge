MODULE := github.com/nmakro/platform2.0-go-challenge

GO111MODULE := on
export GO111MODULE
CGO_ENABLED := 0
export CGO_ENABLED

BINARY_NAME=gwi-server

app:
	@go build -o ${BINARY_NAME} ./cmd/gwiapp

tests:
	@go test --tags=integration_test github.com/nmakro/platform2.0-go-challenge/internal/repositories/maprepo -coverprofile=coverage.out

run:
	@./${BINARY_NAME}

serve:	app run

docker-build:
	@docker build --no-cache -t gwi:latest .

docker-run:
	@docker-compose up --force-recreate

docker-down:
	@docker-compose down --remove-orphans

clean:
	@go clean
	@rm ${BINARY_NAME}
