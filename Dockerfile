FROM golang:1.24 AS builder

COPY . /src
WORKDIR /src

RUN make build

FROM debian:stable-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
		ca-certificates  \
        netbase \
        && rm -rf /var/lib/apt/lists/ \
        && apt-get autoremove -y && apt-get autoclean -y

COPY --from=builder /src/bin/server /app/server

WORKDIR /app

EXPOSE 8000
EXPOSE 9000

CMD ["/app/server", "-conf", "/data/conf"]
