FROM golang:alpine
COPY . /opt/src/shulker
WORKDIR /opt/src/shulker
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o /usr/bin/shulker cmd/main/main.go

FROM busybox:musl
COPY --from=0 /usr/bin/shulker /usr/bin/shulker
VOLUME /opt/shulker
WORKDIR /opt/shulker
CMD ["/usr/bin/shulker", "-config", "/opt/shulker/config.json"]