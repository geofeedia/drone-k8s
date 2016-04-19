## drone-k8s plugin

Plugin for Drone CI to be used in the `publish/deploy` steps.

To run this you will need to clone the repo, build the image, and push it to your associated registry and then use that
as the image for the plugin as shown below. You would replace `image: <your-repo>/<your-org>/drone-k8s:1.0.0` with the appropriate values.

This plugin will perform a rolling update of a pod in a [kubernetes](http://kubernetes.io/) cluster. 

This plugin also assumes your drone server is running inside of [kubernetes](http://kubernetes.io/).

## Available Plugin Options

```no-highlight
replication_controller  -- the name of the rc
docker_image            -- the image name with appropriate repo.
namespace               -- the k8s namespace (defaults to using `default`)
k8s_service_host        -- the K8S_SERVICE_HOST env var
k8s_service_port        -- the K8S_SERVICE_PORT env var
protocol                -- https || http
path_to_cert_authority  -- absolute path to the cert authority (ca.pem)
path_to_client_key      -- absolute path to the client key (worker-key.pem)
path_to_client_cert     -- absolute path to the client cert (worker.pem)
update_period           -- the update period for the rolling update
timeout                 -- the timeout threshold for the rolling update
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
