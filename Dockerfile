FROM debian:bookworm AS main
ENV HOST_CONFIG_DIR=configs
# RUN wget https://github.com/wkhtmltopdf/packaging/releases/download/0.12.6.1-2/wkhtmltox_0.12.6.1-2.bullseye_amd64.deb \
#     && apt-get update \
#     && apt-get install -fy ./wkhtmltox_0.12.6.1-2.bullseye_amd64.deb \
#     && rm -rf ./wkhtmltox_0.12.6.1-2.bullseye_amd64.deb \
#     && apt-get clean -y \
#     && rm -rf /var/lib/apt/lists/*
RUN apt-get update \
    && apt install wget -y \
    && wget https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb \
    && apt-get install -y ./google-chrome-stable_current_amd64.deb \
    && rm -rf google-chrome-stable_current_amd64.deb \
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
RUN go mod download
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
