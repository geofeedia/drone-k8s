package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-plugin-go/plugin"
)

type PluginParams struct {
	ReplicationController  string `json:"replication_controller"`
	DockerImage            string `json:"docker_image"`
	Namespace              string `json:"namespace"`
	K8sServiceHost         string `json:"k8s_service_host"`
	K8sServicePort         string `json:"k8s_service_port"`
	Protocol               string `json:"protocol"`
	KubeConfig             string `json:"kubeconfig"`
	PathToCertAuth         string `json:"path_to_cert_authority"`
	PathToClientKey        string `json:"path_to_client_key"`
	PathToClientCert       string `json:"path_to_client_cert"`
	UpdatePeriod           string `json:"update_period"`
	Timeout                string `json:"timeout"`
	IsDeployment           string `json:"is_deployment"`
	ContainerName          string `json:"container_name"`
	DeploymentResourceName string `json:"deployment_resource_name"`
	EsbConfigPath          string `json:"esb_config_path"`
	ConfigMapName          string `json:"config_map_name"`
	ConfigMapKeyName       string `json:"config_map_key_name"`
	ServiceConfigMapPath   string `json:"service_config_map_path"`
}

func main() {
	fmt.Println("DRONE K8S PLUGIN")

	var (
		repo         = new(drone.Repo)
		build        = new(drone.Build)
		sys          = new(drone.System)
		pluginParams = new(PluginParams)
		cmd          = new(exec.Cmd)
		err          = errors.New("err")
		errMessage   string
		args         = make([]string, 0)
	)

	plugin.Param("build", build)
	plugin.Param("repo", repo)
	plugin.Param("system", sys)
	plugin.Param("vargs", pluginParams)
	plugin.MustParse()

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

	if len(pluginParams.DockerImage) == 0 {
		log.Fatal("No image name provided. Unable to continue.")
	}

	if len(pluginParams.EsbConfigPath) != 0 {
		errMessage = "Unable to update config map with name: " + pluginParams.ConfigMapName

		if len(pluginParams.ConfigMapName) != 0 {
			if len(pluginParams.ConfigMapKeyName) == 0 {
				log.Fatal("No config map key name specified. Unable to replace config map. Exiting.")
			}
			// perform a "dry-run" creation of a new ConfigMap so we can get the yaml output,
			// and then pipe it to the replace command.
			// Clever trick from:  http://stackoverflow.com/questions/38216278/update-k8s-configmap-or-secret-without-deleting-the-existing-one
			args = append(kubectl_global_options(pluginParams), []string{
				"create",
				"configmap",
				pluginParams.ConfigMapName,
				"--from-file=" + pluginParams.ConfigMapKeyName + "=" + pluginParams.EsbConfigPath,
				"--dry-run",
				"-o", "yaml",
			}...)

			createConfigMapCmd := exec.Command("/usr/bin/kubectl", args...)
			trace(createConfigMapCmd)

			args = append(kubectl_global_options(pluginParams), []string{
				"replace",
				"-f", "-",
			}...)
			replaceConfigMapCmd := exec.Command("/usr/bin/kubectl", args...)
			trace(replaceConfigMapCmd)
			success := pipe_commands(createConfigMapCmd, replaceConfigMapCmd)
			if success == nil {
				log.Fatal(errMessage)
			}
		}
		// we need to "kubectl replace" a config map to be used by the pods
	} else if len(pluginParams.ServiceConfigMapPath) != 0 {
		errMessage = "Unable to update config map at location: " + pluginParams.ServiceConfigMapPath
		args = append(kubectl_global_options(pluginParams), []string{
			"replace",
			"-f", pluginParams.ServiceConfigMapPath,
		}...)
		replaceConfigMapCmd := exec.Command("/usr/bin/kubectl", args...)
		trace(replaceConfigMapCmd)
		replaceConfigMapErr := replaceConfigMapCmd.Run()
		if replaceConfigMapErr != nil {
			fmt.Printf("%s\n", replaceConfigMapErr)
			log.Fatal(errMessage)
		}
	}

	if _, parseErr := strconv.ParseBool(pluginParams.IsDeployment); parseErr == nil {
		if len(pluginParams.ContainerName) == 0 || len(pluginParams.DeploymentResourceName) == 0 {
			log.Fatal("Either/both the container name or deployment resource name was/were not provided for deployment patch. Unable to continue.")
		}

		errMessage = "Unable to update deployment for resource " + pluginParams.DeploymentResourceName

		args = append(kubectl_global_options(pluginParams), []string{
			"patch",
			"deployment", pluginParams.DeploymentResourceName,
			"-p", `{"spec":{"template":{"spec":{"containers":[{"name":"` + pluginParams.ContainerName + `","image":"` + pluginParams.DockerImage + `"}]}}}}`,
		}...)
		cmd = exec.Command("/usr/bin/kubectl", args...)
	} else {
		// if not deployment then must be a replication controller
		if len(pluginParams.ReplicationController) == 0 {
			log.Fatal("No replication controller name provided. Unable to continue.")
		}

		errMessage = "Unable to complete rolling-update for " + pluginParams.ReplicationController

		if len(pluginParams.ContainerName) == 0 {
			args = append(kubectl_global_options(pluginParams), []string{
				"rolling-update", pluginParams.ReplicationController,
				"--update-period", pluginParams.UpdatePeriod,
				"--timeout", pluginParams.Timeout,
				"--image", pluginParams.DockerImage,
			}...)
			cmd = exec.Command("/usr/bin/kubectl", args...)
		} else {
			args = append(kubectl_global_options(pluginParams), []string{
				"rolling-update", pluginParams.ReplicationController,
				"--update-period", pluginParams.UpdatePeriod,
				"--timeout", pluginParams.Timeout,
				"--container", pluginParams.ContainerName,
				"--image", pluginParams.DockerImage,
			}...)
			cmd = exec.Command("/usr/bin/kubectl", args...)
		}
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	trace(cmd)
	err = cmd.Run()
	if err != nil {
		fmt.Printf("%s\n", err)
		log.Fatal(errMessage)
	}
}

// Trace writes each command to standard error (preceded by a ‘$ ’) before it
// is executed. Used for debugging the build.
func trace(cmd *exec.Cmd) {
	fmt.Println("$", strings.Join(cmd.Args, " "))
}

// helper to pipe output from one command to the next
func pipe_commands(commands ...*exec.Cmd) []byte {
	for i, command := range commands[:len(commands)-1] {
		out, err := command.StdoutPipe()
		if err != nil {
			return nil
		}
		command.Start()
		commands[i+1].Stdin = out
	}
	final, err := commands[len(commands)-1].Output()
	if err != nil {
		return nil
	}
	return final
}

// helper to add global command options for kubectl
func kubectl_global_options(pluginParams *PluginParams) []string {
	opts := []string{
		"--namespace", pluginParams.Namespace,
		"--server", pluginParams.Protocol + pluginParams.K8sServiceHost + ":" + pluginParams.K8sServicePort,
	}

	if pluginParams.KubeConfig != "" {
		opts = append(opts, "--kubeconfig", pluginParams.KubeConfig)
	}
	if pluginParams.PathToCertAuth != "" {
		opts = append(opts, "--kubeconfig", pluginParams.PathToCertAuth)
	}
	if pluginParams.PathToClientKey != "" {
		opts = append(opts, "--kubeconfig", pluginParams.PathToClientKey)
	}
	if pluginParams.PathToClientCert != "" {
		opts = append(opts, "--kubeconfig", pluginParams.PathToClientCert)
	}

	return opts
}
