# drone-kubernetes
Kubernetes plugin for publishing kubernetes artifacts from [drone](https://drone.io/)

## Overview

This plugin is responsible for publishing artifacts to a kubernetes cluster:

```
sh ./drone-kubernetes <<EOF
{
    "vargs": {
        "replicationcontrollers": [ "example/nginx.json" ],
        "services": [],
        "secrets": ["examples/secrets.yaml"]
        "apiserver": "https://127.0.0.1",
        "namespace": "default",
        "token": "eyJhbGciOiJSUz...",
        "webhook": "https://webhook-gateway.test/drone-deploys
        "webhook_token": "12345abcdf"

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
docker run -i -v $(pwd):/drone/src quay.io/ukhomeofficedigital/drone-kubernetes <<EOF
{
    "vargs": {
        "replicationcontrollers": [ "example/nginx.json" ],
        "services": [],
        "apiserver": "https://127.0.0.1",
        "namespace": "default",
        "token": "eyJhbGciOiJSUz..."
        "webhook": "https://webhook-gateway.test/drone-deploys
        "webhook_token": "12345abcdf"
    }
}
EOF
```

In your `.drone.yml` you will need to add the following

```
deploy:
  kubernetes:
    image: quay.io/ukhomeofficedigital/drone-kubernetes
    replicationcontrollers: ["kubernetes/deep-api-rc.json", "kubernetes/deep-web-rc.json"]
    services: []
    token: $$TOKEN
    apiserver: $$APISERVER
    namespace: default
    webhook: $$WEBHOOK_URL
    webhook_token: $$WEBHOOK_TOKEN
```

The webhook will post a json with the following structure:

```
{
    "Timestamp": 1447753701,
    "Images": [
        "kubernetes/rc.json"
    ],
    "Namespace": "default",
    "Source": "DRONE",
    "Target": "API_SERVER",
    "Url": "WEBOHOOK",
    "Token": "TOKEN"
}
```
