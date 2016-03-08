package main

import(
    "fmt"
    "log"
    "os"
    "os/exec"
    "strings"
    "errors"
    
    "github.com/drone/drone-go/drone"
    "github.com/drone/drone-go/plugin"
)

type PluginParams struct {
    ReplicationController string `json:"replication_controller"`
    DockerImage           string `json:"docker_image"`
    Namespace             string `json:"namespace"`
    K8sServiceHost        string `json:"k8s_service_host"`
    K8sServicePort        string `json:"k8s_service_port"`
    Protocol              string `json:"protocol"`
    PathToCertAuth        string `json:"path_to_cert_authority"`
    PathToClientKey       string `json:"path_to_client_key"`
    PathToClientCert      string `json:"path_to_client_cert"`
    UpdatePeriod          string `json:"update_period"`
    Timeout               string `json:"timeout"`
}

func main() {
    fmt.Println("DRONE K8S PLUGIN")
    
    var (
        workspace    = new(drone.Workspace)
        repo         = new(drone.Repo)
        build        = new(drone.Build)
        sys          = new(drone.System)
        pluginParams = new(PluginParams)
        cmd          = new(exec.Cmd)
        err          = errors.New("err")
    )
    
    plugin.Param("workspace", workspace)
    plugin.Param("build", build)
    plugin.Param("repo", repo)
    plugin.Param("system", sys)
    plugin.Param("vargs", pluginParams)
    plugin.MustParse()
        
    if len(pluginParams.ReplicationController) == 0 {
        log.Fatal("No replication controller name provided. Unable to continue.")
    }
    
    if len(pluginParams.DockerImage) == 0 {
        log.Fatal("No image name provided. Unable to continue.")
    }
    
    if len(pluginParams.Namespace) == 0 {
        pluginParams.Namespace = "default"
    }
    
    if len(pluginParams.UpdatePeriod) == 0 {
        pluginParams.UpdatePeriod = "1m0s"
    }
    
    if len(pluginParams.Timeout) == 0 {
        pluginParams.Timeout = "5m0s"
    }
    
    if len(pluginParams.Protocol) == 0 {
        pluginParams.Protocol = "https://"
    }
    
    if len(pluginParams.K8sServiceHost) == 0 {
        pluginParams.K8sServiceHost = "10.100.0.1"
    }
    
    if len(pluginParams.K8sServicePort) == 0 {
        pluginParams.K8sServicePort = "443"
    }

    cmd = exec.Command(
        "/usr/bin/kubectl",
        "rolling-update", pluginParams.ReplicationController,
        "-s", pluginParams.Protocol + pluginParams.K8sServiceHost + ":" + pluginParams.K8sServicePort,
        "--namespace", pluginParams.Namespace,
        "--certificate-authority", pluginParams.PathToCertAuth,
        "--client-key", pluginParams.PathToClientKey,
        "--client-certificate", pluginParams.PathToClientCert,
        "--update-period", pluginParams.UpdatePeriod,
        "--timeout", pluginParams.Timeout,
        "--image", pluginParams.DockerImage)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    trace(cmd)
    err = cmd.Run()
    if err != nil {
        fmt.Printf("%s", err)
        log.Fatal("Unable to complete kubernetes rolling-update.")
    }
}

// Trace writes each command to standard error (preceded by a ‘$ ’) before it
// is executed. Used for debugging the build.
func trace(cmd *exec.Cmd) {
    fmt.Println("$", strings.Join(cmd.Args, " "))
}