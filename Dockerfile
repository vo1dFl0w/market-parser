FROM golang:1.24.3-alpine AS builder

WORKDIR /market-parser

RUN apk --no-cache add git bash make gcc gettext musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./

ENV CONFIG_PATH=configs/config.yaml
ENV CGO_ENABLED=0

RUN go build --ldflags="-w -s" -o market-parser ./cmd/market-parser

FROM alpine AS runner

RUN apk add --no-cache ca-certificates

WORKDIR /market-parser

COPY --from=builder /market-parser/configs/ /market-parser/configs/
COPY --from=builder /market-parser/market-parser /market-parser/market-parser

EXPOSE 8080

CMD [ "./market-parser" ]