ARG GOLANG_IMAGE="golang:1.22-bookworm"
ARG DISTROLESS_IMAGE="gcr.io/distroless/static-debian12"

FROM ${GOLANG_IMAGE} AS build

WORKDIR /app

ENV CGO_ENABLED=0
ENV GOARCH=amd64
ENV GOOS=linux


COPY ./src/go.mod ./src/go.sum ./
RUN go mod download

COPY ./src/ ./
RUN --mount=type=cache,target="/root/.cache/go-build" go build -o /out/todo .

FROM ${DISTROLESS_IMAGE}

USER 1000
COPY --from=build /out/todo /usr/bin/todo

EXPOSE 8080
ENTRYPOINT ["/usr/bin/todo"]
