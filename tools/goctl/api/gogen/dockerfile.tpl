FROM golang:alpine AS builder

WORKDIR /build

ADD go.mod .
ADD go.sum .
COPY . .
RUN go build -ldflags="-s -w" -o {{.APP_NAME}} .

FROM golang:alpine

WORKDIR /app
COPY --from=builder /build ./
COPY --from=builder /build/etc ./etc

CMD ["./{{.APP_NAME}}"]
