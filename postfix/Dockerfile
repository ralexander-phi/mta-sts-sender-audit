FROM golang:1.23 AS build-stage
WORKDIR /workdir

COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /app


####
FROM python:latest AS template-stage
RUN pip install jinja2-cli
COPY main.cf.tmpl /main.cf.tmpl
ARG name
ARG certdomain
ARG tlsSecurityLevel
RUN jinja2 /main.cf.tmpl                     \
    -D hostname=mail-$name.audit.alexsci.com \
    -D destname=$name.audit.alexsci.com      \
    -D certdomain=$certdomain                \
    -D tlsSecurityLevel=$tlsSecurityLevel    \
    -o /main.cf


####
FROM ubuntu:24.04

RUN --mount=target=/var/lib/apt/lists,type=cache,sharing=locked \
    --mount=target=/var/cache/apt,type=cache,sharing=locked \
    rm -f /etc/apt/apt.conf.d/docker-clean \
    && apt-get update \
    && apt-get -y install postfix

RUN useradd -m -s /bin/bash catchall

COPY entrypoint.sh /entrypoint.sh
COPY --from=template-stage main.cf /etc/postfix/main.cf
COPY --from=build-stage /app /app

EXPOSE 25/tcp 

# Save the Postgres password to somewhere the mailbox_command can access it

CMD [ "/entrypoint.sh" ]

