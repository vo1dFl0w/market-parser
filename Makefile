.PHONY: build run swaggui install-tools ogen

# build app
build:
	go build -o market-parser ./cmd/market-parser

# build+run app
run: build
	SERVER_HTTP_ADDR="localhost:8080" CONFIG_PATH=./configs/config.yaml ./market-parser

# run swagger ui to make requests
swaggui:
	docker run --rm -p 8081:8080 -e SWAGGER_JSON=/openapi.yaml -v ./api/v1/openapi.yaml:/openapi.yaml swaggerapi/swagger-ui

# installing all necessary tools
install-tools:
	go install -v github.com/ogen-go/ogen/cmd/ogen@latest

# generate api
ogen:
	ogen --target ./internal/transport/http/httpgen --package httpgen --clean ./api/v1/openapi.yaml