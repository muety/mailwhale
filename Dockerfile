# Build Stage

FROM golang:1.21 AS api-build-env

WORKDIR /src
ADD . .
RUN CGO_ENABLED=0 go build -o mailwhale

WORKDIR /app
RUN cp /src/mailwhale . && \
    cp /src/version.txt . && \
    cp -r /src/templates .

FROM node:18 AS ui-build-env

WORKDIR /src
ADD webui .
RUN yarn && \
    yarn build

# Run Stage

# When running the application using `docker run`, you can pass environment variables
# to override config values using `-e` syntax.

FROM alpine
WORKDIR /app

ADD config.yml .
COPY --from=api-build-env /app .
COPY --from=ui-build-env /src/public ./webui/public

ENV MW_ENV=prod
ENV MW_WEB_LISTEN_V4=0.0.0.0:3000
ENV MW_STORE_PATH=/data/data.json.db

VOLUME /data

ENTRYPOINT ./mailwhale
