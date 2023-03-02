FROM golang:alpine as builder

WORKDIR /app

COPY . .

RUN go mod tidy

WORKDIR /app/gateway

RUN go build -o gateway-rest

FROM gcr.io/distroless/base-debian11

WORKDIR /

COPY --from=builder /app/gateway/gateway-rest .

ENTRYPOINT [ "/gateway-rest" ]
