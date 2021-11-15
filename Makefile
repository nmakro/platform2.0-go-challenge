MODULE := github.com/nmakro/platform2.0-go-challenge

GO111MODULE := on
export GO111MODULE
CGO_ENABLED := 0
export CGO_ENABLED

BINARY_NAME=app

app:
	go build -o ${BINARY_NAME} ./cmd/wgiapp

run:
	./${BINARY_NAME}

serve:	app run

clean:
	go clean
	rm ${BINARY_NAME}
