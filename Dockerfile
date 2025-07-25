FROM golang:latest AS builder
RUN --mount=type=cache,target=/var/cache/oapi-codegen \
    go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
RUN --mount=type=cache,target=/var/cache/mockery \
    go install github.com/vektra/mockery/v3@latest
RUN --mount=type=cache,target=/var/cache/task \
    go install github.com/go-task/task/v3/cmd/task@latest
WORKDIR /app
COPY ./ ./
RUN task build-docker

FROM ubuntu AS dev
RUN apt-get update && apt-get install -y postgresql-client
WORKDIR /app
COPY --from=builder /app/migrations/ migrations/
COPY --from=builder /app/scripts/wait-for-postgres.sh /app/config/app.yaml /app/bin/jelly ./
EXPOSE 8080 6060
ENTRYPOINT ["./jelly"]

FROM alpine:latest AS prod
WORKDIR /app
COPY --from=builder /app/config/app.yaml config/app.yaml
COPY --from=builder /app/templates/ ./
COPY --from=builder /app/bin/jelly jelly
EXPOSE 8080
ENTRYPOINT ["./jelly"]
