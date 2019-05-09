FROM golang:alpine
RUN apk update && apk add git
COPY . /build
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -a -installsuffix cgo -ldflags='-w -s' -o /build/pgsql-backup-s3

FROM alpine
RUN apk update && apk add postgresql-client
COPY --from=0 /build/pgsql-backup-s3 /usr/bin/pgsql-backup-s3
ENTRYPOINT [ "/usr/bin/pgsql-backup-s3" ]
