# Build the application from source
FROM golang:1.26 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN make release

# Run the tests in the container
FROM build-stage AS test
RUN go test -v ./...

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian13

WORKDIR /

COPY --from=build /app/bin/* /bin/

USER nonroot:nonroot

ENTRYPOINT ["/bin/gowol-server"]
