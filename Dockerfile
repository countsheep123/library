FROM node:12.0.0-alpine as build-frontend

RUN apk add --update --no-cache yarn

WORKDIR /app
COPY ./frontend ./
RUN yarn install
RUN yarn build

# ---

FROM golang:1.12.4-alpine3.9 as build-backend

ENV CGO_ENABLED 0

RUN apk add --update --no-cache \
    curl \
    git \
    openssh-client \
    ca-certificates

WORKDIR /opt
COPY ./backend ./
RUN go mod download
RUN go build -o /server /opt/cmd/server/*.go

ENV MIGRATE_VERSION=v4.3.1
ENV MIGRATE_PLATFORM=linux
ENV MIGRATE_ARCH=amd64
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/${MIGRATE_VERSION}/migrate.${MIGRATE_PLATFORM}-${MIGRATE_ARCH}.tar.gz | tar xvz && \
    mv /opt/migrate.${MIGRATE_PLATFORM}-${MIGRATE_ARCH} /bin/migrate

# ---

FROM scratch

COPY --from=build-frontend /app/dist /static
COPY --from=build-backend /server /server
COPY --from=build-backend /etc/ssl/certs /etc/ssl/certs
COPY --from=build-backend /bin/migrate /bin/migrate

COPY ./migration /opt/migration

CMD ["/server"]
