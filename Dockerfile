FROM golang:alpine AS builder

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
  CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64 \
  GOPROXY=https://goproxy.cn,direct

# Move to working directory /build
WORKDIR /build

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

# Build the application
RUN go build -o main .

# Move to /dist directory as the place for resulting binary folder
WORKDIR /dist

# Copy binary and config file from build to main folder
RUN cp /build/main /build/config.yaml .

# Build a small image
FROM alpine:latest

RUN apk --no-cache add ca-certificates

# Timezone setting
ENV TZ Asia/Shanghai

WORKDIR /root

COPY --from=builder /dist/main /dist/config.yaml ./

# Command to run
CMD ["./main"]