# Docker HUB version control
ARG DEBIAN_VERSION=stable-slim
ARG GOLANG_VERSION=1.24-alpine3.22
ARG ALPINE_VERSION=3.22

FROM debian:${DEBIAN_VERSION} AS asset_builder

# Tooling version control
# TODO: check best practices to use standalone tailwind in production
ENV TAILWINDCSS_VERSION=v4.1.10
ENV DAISYUI_VERSION=v5.0.43
ENV HTMX_VERSION=2.0.6

RUN apt-get update && apt-get install -y curl

WORKDIR /assets

COPY assets/css/input.css ./input.css
COPY templates ./templates

RUN curl -sLo tailwindcss https://github.com/tailwindlabs/tailwindcss/releases/download/${TAILWINDCSS_VERSION}/tailwindcss-linux-x64 && \
    chmod +x ./tailwindcss && \
    curl -sLO https://github.com/saadeghi/daisyui/releases/download/${DAISYUI_VERSION}/daisyui.js && \
    curl -sLO https://github.com/saadeghi/daisyui/releases/download/${DAISYUI_VERSION}/daisyui-theme.js && \
    curl -sLO https://cdn.jsdelivr.net/npm/htmx.org@${HTMX_VERSION}/dist/htmx.min.js

RUN ./tailwindcss -i ./input.css -o ./output.css -m

FROM golang:${GOLANG_VERSION} AS go_builder

# Tooling version control
ENV TEMPL_VERSION=v0.3.898

WORKDIR /app

RUN go install github.com/a-h/templ/cmd/templ@${TEMPL_VERSION}

COPY go.mod go.sum ./
RUN go mod download

# Come from docker-compose, fallbacks to "dev"
ARG ASSETS_VERSION=dev

COPY . .
RUN templ generate
RUN CGO_ENABLED=0 go build -ldflags="-s -w -X main.assetsVersion=${ASSETS_VERSION}" -o ./server ./cmd/web

FROM alpine:${ALPINE_VERSION} AS final_image

WORKDIR /app

COPY --from=go_builder /app/server ./
COPY --from=asset_builder /assets/output.css ./assets/css/
COPY --from=asset_builder /assets/htmx.min.js ./assets/js/
COPY assets/images ./assets/images

ENTRYPOINT ["./server"]
