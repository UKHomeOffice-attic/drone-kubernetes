Use the Kubernetes plugin to deploy pods, replication controllers or services to a kubernetes cluster.
The following parameters are used to configure this plugin:

* `artifacts` - list of artifacts you want to deploy
* `apiserver` - url of kubernetes master
* `namespace` - namespace where the artifacts will be deployed
* `token` - token to validate credentials

The following is a sample Kubernetes configuration in your .drone.yml file:

```yaml
publish:
  kubernetes:
    replicationcontrollers:
        - kubernetes/nginx-rc.json
        - kubernetes/nginx-svc.json
    services:
    apiserver: https://127.0.0.1
    namespace: default
    token: $$TOKEN
```

The plugin assumes the kubernetes artifacts are stored in the current repo.
