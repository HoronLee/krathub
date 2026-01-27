# Build stage
FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS builder

ARG TARGETOS=linux
ARG TARGETARCH

RUN apk add --no-cache make git

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags="-s -w" -o /src/bin/krathub ./app/krathub/service/cmd/server

# Runtime stage
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /src/bin/krathub /app/krathub

# HTTP and gRPC ports
EXPOSE 8000 8001

VOLUME /app/configs

ENV TZ=Asia/Shanghai

CMD ["./krathub", "-conf", "/app/configs"]
