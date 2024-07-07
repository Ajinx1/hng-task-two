# syntax=docker/dockerfile:1

# Create a stage for building the application.
ARG GO_VERSION=1.19.5
FROM golang:${GO_VERSION} AS build
WORKDIR /src

# Copy go.mod and go.sum to the container and download dependencies.
COPY go.mod go.sum ./
RUN go mod download

# Copy the local package files to the container's workspace.
COPY . .

# Build the application.
RUN CGO_ENABLED=0 go build -o /bin/server .


FROM alpine:latest AS final

RUN --mount=type=cache,target=/var/cache/apk \
    apk --update add \
        ca-certificates \
        tzdata \
        && \
        update-ca-certificates

ARG UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    appuser
USER appuser

# Copy the executable from the "build" stage.
COPY --from=build /bin/server /bin/

# Expose the port that the application listens on.
EXPOSE 9090

# What the container should run when it is started.
ENTRYPOINT [ "/bin/server" ]
