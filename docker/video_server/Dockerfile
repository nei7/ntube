FROM golang:1.19.5-alpine3.17 AS builder

RUN apk --no-cache add \
    alpine-sdk \
    librdkafka-dev \
    pkgconf && \
    rm -rf /var/cache/apk/*

WORKDIR /build/

COPY . .

RUN go mod download

RUN GOOS=linux go build -a -ldflags "-extldflags -static" -tags musl \
    github.com/nei7/ntube/cmd/video_server


FROM alpine:3.17 AS certificates

RUN apk --no-cache add ca-certificates


FROM scratch

WORKDIR /api/
ARG VIDEO_STORAGE_PATH

ENV PATH=/api/bin/:$PATH

COPY --from=builder /build/video_server ./bin/video_server

EXPOSE 3001

CMD ["video_server", "-env", "/api/cmd/video_server/.env"]
