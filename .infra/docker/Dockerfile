FROM golang:1.15.8-alpine3.13 as builder

# install utils
RUN apk --update add upx

# copy source files into the container
WORKDIR /go/src
COPY . /go/src

# build app binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/app \
    && upx /go/bin/app


FROM scratch

WORKDIR "/app"

# copy binary and static files from the build container
COPY --from=builder /go/bin/app /app/app
COPY --from=builder /go/src/config.yml /app
COPY --from=builder /go/src/public /app/public
COPY --from=builder /go/src/template /app/template

ENTRYPOINT ["/app/app"]
