FROM golang:1.24-bookworm AS deps

WORKDIR /code

COPY ./go.mod ./go.sum ./

RUN go mod download

FROM golang:1.24-bookworm AS builder

WORKDIR /code
COPY --from=deps /go/pkg /go/pkg

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o ./build/bot_server ./cmd/bot/.

FROM debian:bookworm-slim

WORKDIR /code

RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates && \
    rm -rf /var/lib/apt/lists/*


# Create a non-root user and group
RUN groupadd -r appuser && useradd -r -g appuser appuser


COPY --from=builder /code/build/bot_server .

# Change ownership of the application binary
RUN chown appuser:appuser ./bot_server

USER appuser

EXPOSE 4444

CMD [ "./bot_server" ]
