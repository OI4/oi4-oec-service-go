FROM --platform=$BUILDPLATFORM golang:1.23.0-alpine AS builder

ARG TARGETARCH

WORKDIR /app
COPY . ./
RUN go mod download
RUN go mod verify
RUN GOOS=linux GOARCH=$TARGETARCH go build -a -o output/main demo/main.go

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/output/main .

# Run
ENTRYPOINT ["/app/main"]

