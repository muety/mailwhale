# Build Stage

FROM golang:1.18 AS api-build-env

WORKDIR /src
ADD . .
RUN CGO_ENABLED=0 go build -o mailwhale

WORKDIR /app
RUN cp /src/mailwhale . && \
    cp /src/version.txt . && \
    cp -r /src/templates .

FROM node:14 AS ui-build-env

WORKDIR /src
ADD webui .
RUN yarn && \
    yarn build

# Run Stage

# When running the application using `docker run`, you can pass environment variables
# to override config values using `-e` syntax.

FROM alpine
WORKDIR /app

ENV MW_ENV=prod
ENV MW_SMTP_HOST=''
ENV MW_SMTP_PORT=''
ENV MW_SMTP_USER=''
ENV MW_SMTP_PASS=''
ENV MW_SMTP_TLS=false
ENV MW_WEB_LISTEN_V4=0.0.0.0:3000
ENV MW_WEB_PUBLIC_URL=http://localhost:3000
ENV MW_SECURITY_PEPPER=''
ENV MW_SECURITY_ALLOW_SIGNUP=true
ENV MW_SECURITY_VERIFY_USERS=true
ENV MW_SECURITY_VERIFY_SENDERS=true
ENV MW_STORE_PATH=/data/data.json.db

ADD config.yml .
COPY --from=api-build-env /app .
COPY --from=ui-build-env /src/public ./webui/public

VOLUME /data

ENTRYPOINT ./mailwhale
