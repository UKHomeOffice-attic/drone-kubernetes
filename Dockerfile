# Docker image for Drone's webhook notification plugin
#
#     CGO_ENABLED=0 go build -a -tags netgo
#     docker build --rm=true -t plugins/drone-kubernetes .

FROM gliderlabs/alpine:3.1
RUN apk-install ca-certificates
ADD drone-kubernetes /bin/
ENTRYPOINT ["/bin/drone-kubernetes"]
