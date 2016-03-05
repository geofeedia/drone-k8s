package main

import(
    "fmt"
    "time"
    "log"
    "os"
    "os/exec"
    "io/ioutil"
    "strings"
    
    "github.com/drone/drone-go/drone"
    "github.com/drone/drone-go/plugin"
)

type PluginParams struct {
    ReplicationController string `json:"replication_controller"`
    Registry              string `json:"registry"`
    Image                 string `json:"image"`
    Username              string `json:"username"`
    Email                 string `json:"email"`
    Password              string `json:"password"`
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
    
    if len(pluginParams.Image) == 0 {
        log.Fatal("No image name provided. Unable to continue.")
    }
    
    if len(pluginParams.Registry) == 0 {
        // default docker hub registry if none provided
        pluginParams.Registry = "https://index.docker.io/v1/"
    }
    
    if len(pluginParams.UpdatePeriod) == 0 {
        // default for kubectl rolling-update `--update-period` is 60 seconds
        pluginParams.UpdatePeriod = "1m0s"
    }
    
    if len(pluginParams.Timeout) == 0 {
        // default for kubectl rolling-update `--timeout` is 5 minutes
        pluginParams.Timeout = "5m0s"
    }
    
    if len(pluginParams.Email) == 0 {
        pluginParams.Email = "mail@mail.com"
    }
    
    // ping Docker until available
    for i := 0; i < 3; i++ {
        cmd := exec.Command("/usr/bin/docker", "info")
        cmd.Stdout = ioutil.Discard
        cmd.Stderr = ioutil.Discard
        
        trace(cmd)
        err := cmd.Run()
        if err == nil {
            break
        }
        time.Sleep(time.Second * 5)
    }
    
    if len(pluginParams.Username) != 0 {
        cmd := exec.Command("/usr/bin/docker", "login",
             "-u", pluginParams.Username,
             "-p", pluginParams.Password,
             "-e", pluginParams.Email,
             pluginParams.Registry)
         
        cmd.Stdout = os.Stdout
        cmd.Stderr = os.Stderr
        trace(cmd)
        err := cmd.Run()
        if err != nil {
            fmt.Printf("%s", err)
            log.Fatal("Unable to login to provided Docker registry with given credentials.")
        }
    } else {
        fmt.Println("No username provided so proceeding to use anonymous publishing.")
    }
    
    // kubectl rolling-update --update-period=5s bouncer-e1b-v1.0.0 --image=gcr.io/geofeedia-qa1/service-enterprise-permissions:1.0.0
    cmd := exec.Command(
        "/usr/bin/kubectl",
        "rolling-update", pluginParams.Image,
        "--update-period", pluginParams.UpdatePeriod,
        "--timeout", pluginParams.Timeout,
        pluginParams.ReplicationController,
        "--image", pluginParams.Image)
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