# drone-webhook
Kubernetes plugin for publishing kubernetes artifacts by [@ipedrazas](https://github.com/ipedrazas).

## Overview

This plugin is responsible for publishing artifacts to a kubernetes cluster:

```
sh ./drone-kubernetes <<EOF
{
    "vargs": {
        "replicationcontrollers": [ "example/nginx.json" ],
        "services": [],
        "apiserver": "https://127.0.0.1",
        "namespace": "default",
        "token": "eyJhbGciOiJSUz..."
    }
}
EOF
```

## Docker

Build the Docker container. Note that we need to use the `-netgo` tag so that
the binary is built without a CGO dependency:

```sh
CGO_ENABLED=0 go build -a -tags netgo
docker build --rm=true -t plugins/drone-kubernetes .
```

Deploy to kubernetes:

```
docker run -i -v $(pwd):/drone/src ipedrazas/drone-kubernetes <<EOF
{
    "vargs": {
        "replicationcontrollers": [ "example/nginx.json" ],
        "services": [],
        "apiserver": "https://127.0.0.1",
        "namespace": "default",
        "token": "eyJhbGciOiJSUz..."
    }
}
EOF
```
