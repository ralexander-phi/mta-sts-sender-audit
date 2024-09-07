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
RUN jinja2 /main.cf.tmpl                     \
    -D hostname=mail-$name.audit.alexsci.com \
    -D destname=$name.audit.alexsci.com      \
    -D certdomain=$certdomain                \
    -o /main.cf


####
FROM ubuntu:24.04

RUN --mount=target=/var/lib/apt/lists,type=cache,sharing=locked \
    --mount=target=/var/cache/apt,type=cache,sharing=locked \
    rm -f /etc/apt/apt.conf.d/docker-clean \
    && apt-get update \
    && apt-get -y install postfix

RUN useradd -m -s /bin/bash catchall

COPY --from=template-stage main.cf /etc/postfix/main.cf
COPY --from=build-stage /app /app

EXPOSE 25/tcp 
CMD [ "/usr/sbin/postfix", "-c", "/etc/postfix", "-vv", "start-fg" ]

