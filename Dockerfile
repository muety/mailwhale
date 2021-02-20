# Build Stage

FROM golang:1.15 AS build-env

WORKDIR /src
ADD . .
RUN CGO_ENABLED=0 go build -o mailwhale

WORKDIR /app
RUN cp /src/mailwhale . && \
    cp /src/version.txt .

# Run Stage

# When running the application using `docker run`, you can pass environment variables
# to override config values using `-e` syntax.

FROM alpine
WORKDIR /app

ENV MW_ENV prod
ENV MW_SMTP_HOST ''
ENV MW_SMTP_PORT ''
ENV MW_SMTP_USER ''
ENV MW_SMTP_PASS ''
ENV MW_SMTP_TLS 'false'
ENV MW_WEB_LISTEN_V4 '0.0.0.0:3000'
ENV MW_SECURITY_PEPPER ''
ENV MW_STORE_PATH '/data/data.gob.db'

COPY --from=build-env /app .

VOLUME /data

ENTRYPOINT ./mailwhale
