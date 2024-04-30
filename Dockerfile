FROM golang:1.22.2-alpine3.19 as builder

RUN apk add --no-cache build-base curl libpostal libpostal-dev \
    && /usr/bin/libpostal_data download all /usr/share/libpostal

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=1 go build -o postal-service

FROM alpine:3.19

RUN apk add --no-cache libpostal

WORKDIR /app

EXPOSE 9876

ENTRYPOINT [ "/app/postal-service" ]

COPY --from=builder /usr/share/libpostal /usr/share/libpostal

COPY --from=builder /app/postal-service ./postal-service

