FROM --platform=$BUILDPLATFORM golang:1.16.5-buster AS builder
ARG TARGETPLATFORM
ARG BUILDPLATFORM

# Get dependencies
WORKDIR /opt/build
COPY go.mod go.sum ./
RUN go mod download

# Build binary
COPY ./cmd cmd/
COPY ./pkg pkg/
RUN GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o app ./cmd/limit

# Deploy to minimal environment
FROM --platform=$TARGETPLATFORM debian:buster-slim

EXPOSE 8080
ENV PORT=8080

RUN adduser --system appuser
USER appuser
COPY --from=builder /opt/build/app .

CMD [ "./app" ]
