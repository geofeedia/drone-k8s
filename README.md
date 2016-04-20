## drone-k8s plugin

Plugin for Drone CI to be used in the `publish/deploy` steps that will perform a kubernetes rolling update of the appropriate replication controller and associated pod(s).

To run this you will need to clone the repo, build the image, push it to your associated registry, and then use that
as the image for the plugin as shown below. The easiest way is to use the provided Makefile. Run `$ make` to build and tag the images. 

To use this plugin yourself you would just need to replace `image: your-repo/your-org/drone-k8s:1.0.0` with the appropriate values.

This plugin will perform a rolling update of a pod in a [kubernetes](http://kubernetes.io/) cluster. 

This plugin also assumes your drone server is running inside of [kubernetes](http://kubernetes.io/).

## Available Plugin Options

```no-highlight
replication_controller  -- the name of the rc -- REQUIRED
docker_image            -- the image name with appropriate repo -- REQUIRED
namespace               -- the k8s namespace (defaults to using `default`)
k8s_service_host        -- the K8S_SERVICE_HOST env var (default is 10.100.0.1)
k8s_service_port        -- the K8S_SERVICE_PORT env var (default is 443)
protocol                -- https:// || http:// (default is https://)
path_to_cert_authority  -- absolute path to the cert authority (ca.pem)
path_to_client_key      -- absolute path to the client key (worker-key.pem)
path_to_client_cert     -- absolute path to the client cert (worker.pem)
update_period           -- the update period for the rolling update (default is 1m0s)
timeout                 -- the timeout threshold for the rolling update (default is 5m0s)
```

### Example
```yaml
publish: 
  drone-k8s:
    image: your-repo/your-org/drone-k8s:1.0.0
    replication_controller: some-rc
    namespace: some-ns
    docker_image: some-repo/some-org/some-image:1.0.0
    path_to_cert_authority: /path/to/ca.pem
    path_to_client_key: /path/to/worker-key.pem
    path_to_client_cert: /path/to/worker.pem
    update_period: 5s
    timeout: 30s
```
