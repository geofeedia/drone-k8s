## drone-k8s plugin

Plugin for Drone CI to be used in the `publish/deploy` steps that will perform a kubernetes rolling update of the appropriate replication controller and associated pod(s).

To run this you will need to clone the repo, build the image, push it to your associated registry, and then use that
as the image for the plugin as shown below. The easiest way is to use the provided Makefile. Run `$ make` to build and tag the images. 

To use this plugin yourself you would just need to replace `image: your-repo/your-org/drone-k8s:1.0.0` with the appropriate values.

This plugin will perform a rolling-update for replication controllers or patch of a deployment in a [kubernetes](http://kubernetes.io/) cluster.

To use the plugin you will need to clone the repo, build the image, push it to your associated registry, and then use that
as the image for the plugin as shown below. The easiest way is to use the provided Makefile. Run `$ make` to build and tag the images. 

To use this plugin yourself you would just need to replace `image: your-repo/your-org/drone-k8s:1.0.0` with the appropriate values.

If updating a kubernetes deployment you can specify with `is_deployment`. The default is to perform a rolling-update of a replication controller if `is_deployment` is not set.

This plugin assumes your drone server is running inside of [kubernetes](http://kubernetes.io/).

## Available Plugin Options

```no-highlight
replication_controller   -- REQUIRED for rolling-update: The name of the rc
docker_image             -- REQUIRED for both: The image name with appropriate docker repo
namespace                -- the k8s namespace (defaults to using `default`)
k8s_service_host         -- the K8S_SERVICE_HOST env var (default is 10.100.0.1)
k8s_service_port         -- the K8S_SERVICE_PORT env var (default is 443)
protocol                 -- https || http   (default is https)
path_to_cert_authority   -- absolute path to the cert authority (ca.pem)
path_to_client_key       -- absolute path to the client key (worker-key.pem)
path_to_client_cert      -- absolute path to the client cert (worker.pem)
update_period            -- (only used for rolling-update) the update period for the rolling update (default is 1m0s)
timeout                  -- (only used for rolling-update) the timeout threshold for the rolling update (default is 5m0s)
is_deployment            -- REQUIRED for deployment update: Is this an update of a deployment or not. If not specified then rolling-update of replication controller is assumed.
container_name           -- REQUIRED for deployment update: The name of the container to update the image with
deployment_resource_name -- REQUIRED for deployment update: The name of the deployment resource (i.e. my-deployment.yaml)
```

### Examples

```yaml
# perform a rolling-update
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
    
# perform an update for a deployment
publish: 
  drone-k8s:
    image: your-repo/your-org/drone-k8s:1.0.0
    namespace: some-ns
    is_deployment: true
    deployment_resource_name: some-deployment.yaml
    container_name: some-container
    docker_image: some-repo/some-org/some-image:1.0.0
    path_to_cert_authority: /path/to/ca.pem
    path_to_client_key: /path/to/worker-key.pem
    path_to_client_cert: /path/to/worker.pem
```


Also, if you want to test this locally (replacing the `"vargs"` section with whatever you want to pass in as params for the plugin) you can do...

```
$ go run main.go <<EOF
{
    "repo": {
        "clone_url": "git://github.com/drone/drone",
        "owner": "drone",
        "name": "drone",
        "full_name": "drone/drone"
    },
    "system": {
        "link_url": "https://beta.drone.io"
    },
    "build": {
        "number": 22,
        "status": "success",
        "started_at": 1421029603,
        "finished_at": 1421029813,
        "message": "Update the Readme",
        "author": "johnsmith",
        "author_email": "john.smith@gmail.com",
        "event": "push",
        "branch": "master",
        "commit": "436b7a6e2abaddfd35740527353e78a227ddcb2c",
        "ref": "refs/heads/master"
    },
    "workspace": {
        "root": "/drone/src",
        "path": "/drone/src/github.com/drone/drone"
    },
    "vargs": {
        "is_deployment": "true",
        "deployment_resource_name": "my-resource.yaml",
        "container_name": "containerName",
        "docker_image": "myImage"
    }
}
EOF
```

