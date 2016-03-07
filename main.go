package main

import(
    "fmt"
    "log"
    "os"
    "os/exec"
    "strings"
    
    "github.com/drone/drone-go/drone"
    "github.com/drone/drone-go/plugin"
)

type PluginParams struct {
    ReplicationController string `json:"replication_controller"`
    DockerImage              string `json:"docker_image"`
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
    
    // kubectl rolling-update --update-period=5s YOUR-RC-e1b-v1.0.0 --image=registry/org/repo:1.0.0
    cmd := exec.Command(
        "/usr/bin/kubectl",
        "rolling-update", pluginParams.ReplicationController,
        "--update-period", pluginParams.UpdatePeriod,
        "--timeout", pluginParams.Timeout,
        "--image", pluginParams.DockerImage)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    
    trace(cmd)
    err := cmd.Run()
    if err != nil {
        fmt.Printf("%s", err)
        log.Fatal("Unable to complete kubernetes rolling-update.")
    }
}

// Trace writes each command to standard error (preceded by a ‘$ ’) before it
// is executed. Used for debugging your build.
func trace(cmd *exec.Cmd) {
    fmt.Println("$", strings.Join(cmd.Args, " "))
}