FROM debian:bookworm AS main
ENV HOST_CONFIG_DIR=configs
RUN apt-get update \
    && apt-get install wget -y \
    && wget https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb \
    && apt-get install ./google-chrome-stable_current_amd64.deb -y \
    && rm -rf google-chrome-stable_current_amd64.deb \
    && apt-get clean -y \
    && rm -rf /var/lib/apt/lists/*
RUN wget https://github.com/wkhtmltopdf/packaging/releases/download/0.12.6.1-3/wkhtmltox_0.12.6.1-3.bookworm_amd64.deb \
    && apt-get update \
    && apt-get install -fy ./wkhtmltox_0.12.6.1-3.bookworm_amd64.deb \
    && rm -rf ./wkhtmltox_0.12.6.1-3.bookworm_amd64.deb \
    && apt-get clean -y \
    && rm -rf /var/lib/apt/lists/*

FROM golang:1.22-bookworm AS dev-build
RUN apt update \
    && apt-get install build-essential -y \
    && apt-get clean -y \
    && rm -rf /var/lib/apt/lists/*
WORKDIR /build
COPY go.mod .
COPY go.sum .
COPY Makefile .
RUN make setup
COPY . .

FROM dev-build as dev
COPY --from=main / /
CMD ["go", "run", "./cmd/pdf-service/main.go"]

FROM dev-build AS prod-build
RUN make

FROM main AS prod
WORKDIR /etc/pdf-service
COPY --from=prod-build /build/templates ./templates
COPY --from=prod-build /build/configs ./configs
COPY --from=prod-build /build/bin/pdf-service /usr/local/bin
CMD ["pdf-service"]
