# syntax=docker/dockerfile:1

# Build the application from source
FROM --platform=$BUILDPLATFORM golang:1.25 AS build-stage
ARG TARGETOS
ARG TARGETARCH
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} CGO_ENABLED=0 go build -o /dice

# Run the tests in the container
FROM build-stage AS run-test-stage
RUN go test -v ./...

# Deploy the application binary into a lean image
FROM alpine:3.23 AS build-release-stage

WORKDIR /

COPY --from=build-stage /dice /dice

EXPOSE 8080

# USER nonroot:nonroot

ENTRYPOINT ["/dice"]
