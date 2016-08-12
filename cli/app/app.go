/*
Copyright 2016 Skippbox, Ltd All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package app

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/docker/docker/api/client/bundlefile"
	"github.com/docker/libcompose/config"
	"github.com/docker/libcompose/docker"
	"github.com/docker/libcompose/lookup"
	"github.com/docker/libcompose/project"

	"encoding/json"
	"io/ioutil"

	// install kubernetes api
	_ "k8s.io/kubernetes/pkg/api/install"
	_ "k8s.io/kubernetes/pkg/apis/extensions/install"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/apis/extensions"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
	"k8s.io/kubernetes/pkg/runtime"
	"k8s.io/kubernetes/pkg/util/intstr"

	deployapi "github.com/openshift/origin/pkg/deploy/api"
	// install kubernetes api
	_ "github.com/openshift/origin/pkg/deploy/api/install"

	"github.com/fatih/structs"
	"github.com/ghodss/yaml"
)

const (
	letterBytes        = "abcdefghijklmnopqrstuvwxyz0123456789"
	DefaultComposeFile = "docker-compose.yml"
)

var unsupportedKey = map[string]int{
	"Build":         0,
	"CapAdd":        0,
	"CapDrop":       0,
	"CPUSet":        0,
	"CPUShares":     0,
	"CPUQuota":      0,
	"CgroupParent":  0,
	"Devices":       0,
	"DependsOn":     0,
	"DNS":           0,
	"DNSSearch":     0,
	"DomainName":    0,
	"Entrypoint":    0,
	"EnvFile":       0,
	"Expose":        0,
	"Extends":       0,
	"ExternalLinks": 0,
	"ExtraHosts":    0,
	"Hostname":      0,
	"Ipc":           0,
	"Logging":       0,
	"MacAddress":    0,
	"MemLimit":      0,
	"MemSwapLimit":  0,
	"NetworkMode":   0,
	"Networks":      0,
	"Pid":           0,
	"SecurityOpt":   0,
	"ShmSize":       0,
	"StopSignal":    0,
	"VolumeDriver":  0,
	"VolumesFrom":   0,
	"Uts":           0,
	"ReadOnly":      0,
	"StdinOpen":     0,
	"Tty":           0,
	"User":          0,
	"Ulimits":       0,
	"Dockerfile":    0,
	"Net":           0,
	"Args":          0,
}

var composeOptions = map[string]string{
	"Build":         "build",
	"CapAdd":        "cap_add",
	"CapDrop":       "cap_drop",
	"CPUSet":        "cpuset",
	"CPUShares":     "cpu_shares",
	"CPUQuota":      "cpu_quota",
	"CgroupParent":  "cgroup_parent",
	"Devices":       "devices",
	"DependsOn":     "depends_on",
	"DNS":           "dns",
	"DNSSearch":     "dns_search",
	"DomainName":    "domainname",
	"Entrypoint":    "entrypoint",
	"EnvFile":       "env_file",
	"Expose":        "expose",
	"Extends":       "extends",
	"ExternalLinks": "external_links",
	"ExtraHosts":    "extra_hosts",
	"Hostname":      "hostname",
	"Ipc":           "ipc",
	"Logging":       "logging",
	"MacAddress":    "mac_address",
	"MemLimit":      "mem_limit",
	"MemSwapLimit":  "memswap_limit",
	"NetworkMode":   "network_mode",
	"Networks":      "networks",
	"Pid":           "pid",
	"SecurityOpt":   "security_opt",
	"ShmSize":       "shm_size",
	"StopSignal":    "stop_signal",
	"VolumeDriver":  "volume_driver",
	"VolumesFrom":   "volumes_from",
	"Uts":           "uts",
	"ReadOnly":      "read_only",
	"StdinOpen":     "stdin_open",
	"Tty":           "tty",
	"User":          "user",
	"Ulimits":       "ulimits",
	"Dockerfile":    "dockerfile",
	"Net":           "net",
	"Args":          "args",
}

// RandStringBytes generates randomly n-character string
func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// BeforeApp is an action that is executed before any cli command.
func BeforeApp(c *cli.Context) error {
	if c.GlobalBool("verbose") {
		logrus.SetLevel(logrus.DebugLevel)
	}
	return nil
}

// Ps lists all rc, svc.
func Ps(c *cli.Context) {
	//factory := cmdutil.NewFactory(nil)
	//clientConfig, err := factory.ClientConfig()
	//if err != nil {
	//	logrus.Fatalf("Failed to get Kubernetes client config: %v", err)
	//}
	//client := client.NewOrDie(clientConfig)
	//
	//if c.BoolT("svc") {
	//	fmt.Printf("%-20s%-20s%-20s%-20s\n", "Name", "Cluster IP", "Ports", "Selectors")
	//	for name := range p.Configs {
	//		var ports string
	//		var selectors string
	//		services, err := client.Services(api.NamespaceDefault).Get(name)
	//
	//		if err != nil {
	//			logrus.Debugf("Cannot find service for: ", name)
	//		} else {
	//
	//			for i := range services.Spec.Ports {
	//				p := strconv.Itoa(int(services.Spec.Ports[i].Port))
	//				ports += ports + string(services.Spec.Ports[i].Protocol) + "(" + p + "),"
	//			}
	//
	//			for k, v := range services.ObjectMeta.Labels {
	//				selectors += selectors + k + "=" + v + ","
	//			}
	//
	//			ports = strings.TrimSuffix(ports, ",")
	//			selectors = strings.TrimSuffix(selectors, ",")
	//
	//			fmt.Printf("%-20s%-20s%-20s%-20s\n", services.ObjectMeta.Name,
	//				services.Spec.ClusterIP, ports, selectors)
	//		}
	//
	//	}
	//}
	//
	//if c.BoolT("rc") {
	//	fmt.Printf("%-15s%-15s%-30s%-10s%-20s\n", "Name", "Containers", "Images",
	//		"Replicas", "Selectors")
	//	for name := range p.Configs {
	//		var selectors string
	//		var containers string
	//		var images string
	//		rc, err := client.ReplicationControllers(api.NamespaceDefault).Get(name)
	//
	//		/* Should grab controller, container, image, selector, replicas */
	//
	//		if err != nil {
	//			logrus.Debugf("Cannot find rc for: ", string(name))
	//		} else {
	//
	//			for k, v := range rc.Spec.Selector {
	//				selectors += selectors + k + "=" + v + ","
	//			}
	//
	//			for i := range rc.Spec.Template.Spec.Containers {
	//				c := rc.Spec.Template.Spec.Containers[i]
	//				containers += containers + c.Name + ","
	//				images += images + c.Image + ","
	//			}
	//			selectors = strings.TrimSuffix(selectors, ",")
	//			containers = strings.TrimSuffix(containers, ",")
	//			images = strings.TrimSuffix(images, ",")
	//
	//			fmt.Printf("%-15s%-15s%-30s%-10d%-20s\n", rc.ObjectMeta.Name, containers,
	//				images, rc.Spec.Replicas, selectors)
	//		}
	//	}
	//}

}

// Delete deletes all rc, svc.
func Delete(c *cli.Context) {
	//factory := cmdutil.NewFactory(nil)
	//clientConfig, err := factory.ClientConfig()
	//if err != nil {
	//	logrus.Fatalf("Failed to get Kubernetes client config: %v", err)
	//}
	//client := client.NewOrDie(clientConfig)
	//
	//for name := range p.Configs {
	//	if len(c.String("name")) > 0 && name != c.String("name") {
	//		continue
	//	}
	//
	//	if c.BoolT("svc") {
	//		err := client.Services(api.NamespaceDefault).Delete(name)
	//		if err != nil {
	//			logrus.Fatalf("Unable to delete service %s: %s\n", name, err)
	//		}
	//	} else if c.BoolT("rc") {
	//		err := client.ReplicationControllers(api.NamespaceDefault).Delete(name)
	//		if err != nil {
	//			logrus.Fatalf("Unable to delete replication controller %s: %s\n", name, err)
	//		}
	//	}
	//}
}

// Scale scales rc.
func Scale(c *cli.Context) {
	//factory := cmdutil.NewFactory(nil)
	//clientConfig, err := factory.ClientConfig()
	//if err != nil {
	//	logrus.Fatalf("Failed to get Kubernetes client config: %v", err)
	//}
	//client := client.NewOrDie(clientConfig)
	//
	//if c.Int("scale") <= 0 {
	//	logrus.Fatalf("Scale must be defined and a positive number")
	//}
	//
	//for name := range p.Configs {
	//	if len(c.String("rc")) == 0 || c.String("rc") == name {
	//		s, err := client.ExtensionsClient.Scales(api.NamespaceDefault).Get("ReplicationController", name)
	//		if err != nil {
	//			logrus.Fatalf("Error retrieving scaling data: %s\n", err)
	//		}
	//
	//		s.Spec.Replicas = int32(c.Int("scale"))
	//
	//		s, err = client.ExtensionsClient.Scales(api.NamespaceDefault).Update("ReplicationController", s)
	//		if err != nil {
	//			logrus.Fatalf("Error updating scaling data: %s\n", err)
	//		}
	//
	//		fmt.Printf("Scaling %s to: %d\n", name, s.Spec.Replicas)
	//	}
	//}
}

// Create the file to write to if --out is specified
func createOutFile(out string) *os.File {
	var f *os.File
	var err error
	if len(out) != 0 {
		f, err = os.Create(out)
		if err != nil {
			logrus.Fatalf("error opening file: %v", err)
		}
	}
	return f
}

// Init RC object
func initRC(name string, service ServiceConfig, replicas int) *api.ReplicationController {
	rc := &api.ReplicationController{
		TypeMeta: unversioned.TypeMeta{
			Kind:       "ReplicationController",
			APIVersion: "v1",
		},
		ObjectMeta: api.ObjectMeta{
			Name: name,
			//Labels: map[string]string{"service": name},
		},
		Spec: api.ReplicationControllerSpec{
			Selector: map[string]string{"service": name},
			Replicas: int32(replicas),
			Template: &api.PodTemplateSpec{
				ObjectMeta: api.ObjectMeta{
				//Labels: map[string]string{"service": name},
				},
				Spec: api.PodSpec{
					Containers: []api.Container{
						{
							Name:  name,
							Image: service.Image,
						},
					},
				},
			},
		},
	}
	return rc
}

// Init SC object
func initSC(name string, service ServiceConfig) *api.Service {
	sc := &api.Service{
		TypeMeta: unversioned.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: api.ObjectMeta{
			Name: name,
			//Labels: map[string]string{"service": name},
		},
		Spec: api.ServiceSpec{
			Selector: map[string]string{"service": name},
		},
	}
	return sc
}

// Init DC object
func initDC(name string, service ServiceConfig, replicas int) *extensions.Deployment {
	dc := &extensions.Deployment{
		TypeMeta: unversioned.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "extensions/v1beta1",
		},
		ObjectMeta: api.ObjectMeta{
			Name:   name,
			Labels: map[string]string{"service": name},
		},
		Spec: extensions.DeploymentSpec{
			Replicas: int32(replicas),
			Selector: &unversioned.LabelSelector{
				MatchLabels: map[string]string{"service": name},
			},
			//UniqueLabelKey: p.Name,
			Template: api.PodTemplateSpec{
				ObjectMeta: api.ObjectMeta{
					Labels: map[string]string{"service": name},
				},
				Spec: api.PodSpec{
					Containers: []api.Container{
						{
							Name:  name,
							Image: service.Image,
						},
					},
				},
			},
		},
	}
	return dc
}

// Init DS object
func initDS(name string, service ServiceConfig) *extensions.DaemonSet {
	ds := &extensions.DaemonSet{
		TypeMeta: unversioned.TypeMeta{
			Kind:       "DaemonSet",
			APIVersion: "extensions/v1beta1",
		},
		ObjectMeta: api.ObjectMeta{
			Name: name,
		},
		Spec: extensions.DaemonSetSpec{
			Template: api.PodTemplateSpec{
				ObjectMeta: api.ObjectMeta{
					Name: name,
				},
				Spec: api.PodSpec{
					Containers: []api.Container{
						{
							Name:  name,
							Image: service.Image,
						},
					},
				},
			},
		},
	}
	return ds
}

// Init RS object
func initRS(name string, service ServiceConfig, replicas int) *extensions.ReplicaSet {
	rs := &extensions.ReplicaSet{
		TypeMeta: unversioned.TypeMeta{
			Kind:       "ReplicaSet",
			APIVersion: "extensions/v1beta1",
		},
		ObjectMeta: api.ObjectMeta{
			Name: name,
		},
		Spec: extensions.ReplicaSetSpec{
			Replicas: int32(replicas),
			Selector: &unversioned.LabelSelector{
				MatchLabels: map[string]string{"service": name},
			},
			Template: api.PodTemplateSpec{
				ObjectMeta: api.ObjectMeta{},
				Spec: api.PodSpec{
					Containers: []api.Container{
						{
							Name:  name,
							Image: service.Image,
						},
					},
				},
			},
		},
	}
	return rs
}

// initDeploymentConfig initialize OpenShifts DeploymentConfig object
func initDeploymentConfig(name string, service ServiceConfig, replicas int) *deployapi.DeploymentConfig {
	dc := &deployapi.DeploymentConfig{
		TypeMeta: unversioned.TypeMeta{
			Kind:       "DeploymentConfig",
			APIVersion: "v1",
		},
		ObjectMeta: api.ObjectMeta{
			Name:   name,
			Labels: map[string]string{"service": name},
		},
		Spec: deployapi.DeploymentConfigSpec{
			Replicas: int32(replicas),
			Selector: map[string]string{"service": name},
			//UniqueLabelKey: p.Name,
			Template: &api.PodTemplateSpec{
				ObjectMeta: api.ObjectMeta{
					Labels: map[string]string{"service": name},
				},
				Spec: api.PodSpec{
					Containers: []api.Container{
						{
							Name:  name,
							Image: service.Image,
						},
					},
				},
			},
		},
	}
	return dc
}

// Configure the environment variables.
func configEnvs(name string, service ServiceConfig) []api.EnvVar {
	envs := []api.EnvVar{}
	for _, v := range service.Environment {
		envs = append(envs, api.EnvVar{
			Name:  v.Name,
			Value: v.Value,
		})
	}

	return envs
}

// Configure the container volumes.
func configVolumes(service ServiceConfig) ([]api.VolumeMount, []api.Volume) {
	volumesMount := []api.VolumeMount{}
	volumes := []api.Volume{}
	volumeSource := api.VolumeSource{}
	for _, volume := range service.Volumes {
		name, host, container, mode, err := parseVolume(volume)
		if err != nil {
			logrus.Warningf("Failed to configure container volume: %v", err)
			continue
		}

		// if volume name isn't specified, set it to a random string of 20 chars
		if len(name) == 0 {
			name = RandStringBytes(20)
		}
		// check if ro/rw mode is defined, default rw
		readonly := len(mode) > 0 && mode == "ro"

		volumesMount = append(volumesMount, api.VolumeMount{Name: name, ReadOnly: readonly, MountPath: container})

		if len(host) > 0 {
			volumeSource = api.VolumeSource{HostPath: &api.HostPathVolumeSource{Path: host}}
		} else {
			volumeSource = api.VolumeSource{EmptyDir: &api.EmptyDirVolumeSource{}}
		}

		volumes = append(volumes, api.Volume{Name: name, VolumeSource: volumeSource})
	}
	return volumesMount, volumes
}

// parseVolume parse a given volume, which might be [name:][host:]container[:access_mode]
func parseVolume(volume string) (name, host, container, mode string, err error) {
	separator := ":"
	volumeStrings := strings.Split(volume, separator)
	if len(volumeStrings) == 0 {
		return
	}
	// Set name if existed
	if !isPath(volumeStrings[0]) {
		name = volumeStrings[0]
		volumeStrings = volumeStrings[1:]
	}
	if len(volumeStrings) == 0 {
		err = fmt.Errorf("invalid volume format: %s", volume)
		return
	}
	if volumeStrings[len(volumeStrings)-1] == "rw" || volumeStrings[len(volumeStrings)-1] == "ro" {
		mode = volumeStrings[len(volumeStrings)-1]
		volumeStrings = volumeStrings[:len(volumeStrings)-1]
	}
	container = volumeStrings[len(volumeStrings)-1]
	volumeStrings = volumeStrings[:len(volumeStrings)-1]
	if len(volumeStrings) == 1 {
		host = volumeStrings[0]
	}
	if !isPath(container) || (len(host) > 0 && !isPath(host)) || len(volumeStrings) > 1 {
		err = fmt.Errorf("invalid volume format: %s", volume)
		return
	}
	return
}

func isPath(substring string) bool {
	return strings.Contains(substring, "/")
}

// Configure the container ports.
func configPorts(name string, service ServiceConfig) []api.ContainerPort {
	ports := []api.ContainerPort{}
	for _, port := range service.Port {
		var p api.Protocol
		switch port.Protocol {
		default:
			p = api.ProtocolTCP
		case ProtocolTCP:
			p = api.ProtocolTCP
		case ProtocolUDP:
			p = api.ProtocolUDP
		}
		ports = append(ports, api.ContainerPort{
			ContainerPort: port.ContainerPort,
			Protocol:      p,
		})
	}

	return ports
}

// Configure the container service ports.
func configServicePorts(name string, service ServiceConfig) []api.ServicePort {
	servicePorts := []api.ServicePort{}
	for _, port := range service.Port {
		if port.HostPort == 0 {
			port.HostPort = port.ContainerPort
		}
		var p api.Protocol
		switch port.Protocol {
		default:
			p = api.ProtocolTCP
		case ProtocolTCP:
			p = api.ProtocolTCP
		case ProtocolUDP:
			p = api.ProtocolUDP
		}
		var targetPort intstr.IntOrString
		targetPort.IntVal = port.ContainerPort
		targetPort.StrVal = strconv.Itoa(int(port.ContainerPort))
		servicePorts = append(servicePorts, api.ServicePort{
			Name:       strconv.Itoa(int(port.HostPort)),
			Protocol:   p,
			Port:       port.HostPort,
			TargetPort: targetPort,
		})
	}
	return servicePorts
}

// Transform data to json/yaml
func transformer(obj runtime.Object, generateYaml bool) ([]byte, error) {
	//  Convert to versioned object
	objectVersion := obj.GetObjectKind().GroupVersionKind()
	version := unversioned.GroupVersion{Group: objectVersion.Group, Version: objectVersion.Version}
	versionedObj, err := api.Scheme.ConvertToVersion(obj, version)
	if err != nil {
		return nil, err
	}

	// convert data to json / yaml
	data, err := json.MarshalIndent(versionedObj, "", "  ")
	if generateYaml == true {
		data, err = yaml.Marshal(versionedObj)
	}
	if err != nil {
		return nil, err
	}
	logrus.Debugf("%s\n", data)
	return data, nil
}

// load Environment Variable from bundles file
func loadEnvVars(service bundlefile.Service) ([]EnvVar, string) {
	envs := []EnvVar{}
	for _, env := range service.Env {
		character := "="
		if strings.Contains(env, character) {
			value := env[strings.Index(env, character)+1:]
			name := env[0:strings.Index(env, character)]
			name = strings.TrimSpace(name)
			value = strings.TrimSpace(value)
			envs = append(envs, EnvVar{
				Name:  name,
				Value: value,
			})
		} else {
			character = ":"
			if strings.Contains(env, character) {
				charQuote := "'"
				value := env[strings.Index(env, character)+1:]
				name := env[0:strings.Index(env, character)]
				name = strings.TrimSpace(name)
				value = strings.TrimSpace(value)
				if strings.Contains(value, charQuote) {
					value = strings.Trim(value, "'")
				}
				envs = append(envs, EnvVar{
					Name:  name,
					Value: value,
				})
			} else {
				return envs, "Invalid container env " + env
			}
		}
	}
	return envs, ""
}

// load Environment Variable from compose file
func loadEnvVarsFromCompose(e map[string]string) []EnvVar {
	envs := []EnvVar{}
	for k, v := range e {
		envs = append(envs, EnvVar{
			Name:  k,
			Value: v,
		})
	}
	return envs
}

// load Ports from bundles file
func loadPorts(service bundlefile.Service) ([]Ports, string) {
	ports := []Ports{}
	for _, port := range service.Ports {
		var p Protocol
		switch port.Protocol {
		default:
			p = ProtocolTCP
		case "TCP":
			p = ProtocolTCP
		case "UDP":
			p = ProtocolUDP
		}
		ports = append(ports, Ports{
			HostPort:      int32(port.Port),
			ContainerPort: int32(port.Port),
			Protocol:      p,
		})
	}
	return ports, ""
}

// Load Ports from compose file
func loadPortsFromCompose(composePorts []string) ([]Ports, string) {
	ports := []Ports{}
	character := ":"
	for _, port := range composePorts {
		p := ProtocolTCP
		if strings.Contains(port, character) {
			hostPort := port[0:strings.Index(port, character)]
			hostPort = strings.TrimSpace(hostPort)
			hostPortInt, err := strconv.Atoi(hostPort)
			if err != nil {
				return nil, "Invalid host port of " + port
			}
			containerPort := port[strings.Index(port, character)+1:]
			containerPort = strings.TrimSpace(containerPort)
			containerPortInt, err := strconv.Atoi(containerPort)
			if err != nil {
				return nil, "Invalid container port of " + port
			}
			ports = append(ports, Ports{
				HostPort:      int32(hostPortInt),
				ContainerPort: int32(containerPortInt),
				Protocol:      p,
			})
		} else {
			containerPortInt, err := strconv.Atoi(port)
			if err != nil {
				return nil, "Invalid container port of " + port
			}
			ports = append(ports, Ports{
				ContainerPort: int32(containerPortInt),
				Protocol:      p,
			})
		}

	}
	return ports, ""
}

// load Image from bundles file
func loadImage(service bundlefile.Service) (string, string) {
	character := "@"
	if strings.Contains(service.Image, character) {
		return service.Image[0:strings.Index(service.Image, character)], ""
	}
	return "", "Invalid image format"
}

// Load DAB file into KomposeObject
func loadBundlesFile(file string) KomposeObject {
	komposeObject := KomposeObject{
		ServiceConfigs: make(map[string]ServiceConfig),
	}
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		logrus.Fatalf("Failed to read bundles file: %v", err)
	}
	reader := strings.NewReader(string(buf))
	bundle, err := bundlefile.LoadFile(reader)
	if err != nil {
		logrus.Fatalf("Failed to parse bundles file: %v", err)
	}

	for name, service := range bundle.Services {
		checkUnsupportedKey(service)
		serviceConfig := ServiceConfig{}
		serviceConfig.Command = service.Command
		serviceConfig.Args = service.Args
		// convert bundle labels to annotations
		serviceConfig.Annotations = service.Labels

		image, err := loadImage(service)
		if err != "" {
			logrus.Fatalf("Failed to load image from bundles file: %v", err)
		}
		serviceConfig.Image = image

		envs, err := loadEnvVars(service)
		if err != "" {
			logrus.Fatalf("Failed to load envvar from bundles file: %v", err)
		}
		serviceConfig.Environment = envs

		ports, err := loadPorts(service)
		if err != "" {
			logrus.Fatalf("Failed to load ports from bundles file: %v", err)
		}
		serviceConfig.Port = ports

		if service.WorkingDir != nil {
			serviceConfig.WorkingDir = *service.WorkingDir
		}

		komposeObject.ServiceConfigs[name] = serviceConfig
	}
	return komposeObject
}

// Load compose file into KomposeObject
func loadComposeFile(file string) KomposeObject {
	komposeObject := KomposeObject{
		ServiceConfigs: make(map[string]ServiceConfig),
	}
	context := &docker.Context{}
	if file == "" {
		file = "docker-compose.yml"
	}
	context.ComposeFiles = []string{file}

	if context.ResourceLookup == nil {
		context.ResourceLookup = &lookup.FileResourceLookup{}
	}

	if context.EnvironmentLookup == nil {
		cwd, err := os.Getwd()
		if err != nil {
			return KomposeObject{}
		}
		context.EnvironmentLookup = &lookup.ComposableEnvLookup{
			Lookups: []config.EnvironmentLookup{
				&lookup.EnvfileLookup{
					Path: filepath.Join(cwd, ".env"),
				},
				&lookup.OsEnvLookup{},
			},
		}
	}

	// load compose file into composeObject
	composeObject := project.NewProject(&context.Context, nil, nil)
	err := composeObject.Parse()
	if err != nil {
		logrus.Fatalf("Failed to load compose file: %v", err)
	}

	// transform composeObject into komposeObject
	composeServiceNames := composeObject.ServiceConfigs.Keys()

	// volume config and network config are not supported
	if len(composeObject.NetworkConfigs) > 0 {
		logrus.Warningf("Unsupported network configuration of compose v2 - ignoring")
	}
	if len(composeObject.VolumeConfigs) > 0 {
		logrus.Warningf("Unsupported volume configuration of compose v2 - ignoring")
	}

	networksWarningFound := false
	for _, name := range composeServiceNames {
		if composeServiceConfig, ok := composeObject.ServiceConfigs.Get(name); ok {
			//FIXME: networks always contains one default element, even it isn't declared in compose v2.
			if composeServiceConfig.Networks != nil && len(composeServiceConfig.Networks.Networks) > 0 &&
				composeServiceConfig.Networks.Networks[0].Name != "default" &&
				!networksWarningFound {
				logrus.Warningf("Unsupported key networks - ignoring")
				networksWarningFound = true
			}
			checkUnsupportedKey(composeServiceConfig)
			serviceConfig := ServiceConfig{}
			serviceConfig.Image = composeServiceConfig.Image
			serviceConfig.ContainerName = composeServiceConfig.ContainerName

			// load environment variables
			envs := loadEnvVarsFromCompose(composeServiceConfig.Environment.ToMap())
			serviceConfig.Environment = envs

			// load ports
			ports, err := loadPortsFromCompose(composeServiceConfig.Ports)
			if err != "" {
				logrus.Fatalf("Failed to load ports from compose file: %v", err)
			}
			serviceConfig.Port = ports

			serviceConfig.WorkingDir = composeServiceConfig.WorkingDir
			serviceConfig.Volumes = composeServiceConfig.Volumes

			// convert compose labels to annotations
			serviceConfig.Annotations = map[string]string(composeServiceConfig.Labels)

			serviceConfig.CPUSet = composeServiceConfig.CPUSet
			serviceConfig.CPUShares = composeServiceConfig.CPUShares
			serviceConfig.CPUQuota = composeServiceConfig.CPUQuota
			serviceConfig.CapAdd = composeServiceConfig.CapAdd
			serviceConfig.CapDrop = composeServiceConfig.CapDrop
			serviceConfig.Expose = composeServiceConfig.Expose
			serviceConfig.Privileged = composeServiceConfig.Privileged
			serviceConfig.Restart = composeServiceConfig.Restart
			serviceConfig.User = composeServiceConfig.User

			komposeObject.ServiceConfigs[name] = serviceConfig
		}
	}
	return komposeObject
}

type convertOptions struct {
	toStdout               bool
	createD                bool
	createRC               bool
	createDS               bool
	createDeploymentConfig bool
	createChart            bool
	generateYaml           bool
	replicas               int
	inputFile              string
	outFile                string
}

// Convert komposeObject to K8S controllers
func komposeConvert(komposeObject KomposeObject, opt convertOptions) []runtime.Object {
	var svcnames []string

	// this will hold all the converted data
	var allobjects []runtime.Object
	for name, service := range komposeObject.ServiceConfigs {
		var objects []runtime.Object
		svcnames = append(svcnames, name)
		sc := initSC(name, service)

		if opt.createD {
			objects = append(objects, initDC(name, service, opt.replicas))
		}
		if opt.createDS {
			objects = append(objects, initDS(name, service))
		}
		if opt.createRC {
			objects = append(objects, initRC(name, service, opt.replicas))
		}
		if opt.createDeploymentConfig {
			objects = append(objects, initDeploymentConfig(name, service, opt.replicas)) // OpenShift DeploymentConfigs
		}

		// Configure the environment variables.
		envs := configEnvs(name, service)

		// Configure the container command.
		var cmds []string
		for _, cmd := range service.Command {
			cmds = append(cmds, cmd)
		}
		// Configure the container volumes.
		volumesMount, volumes := configVolumes(service)

		// Configure the container ports.
		ports := configPorts(name, service)

		// Configure the service ports.
		servicePorts := configServicePorts(name, service)
		sc.Spec.Ports = servicePorts

		// Configure labels
		labels := map[string]string{"service": name}
		sc.ObjectMeta.Labels = labels
		// Configure annotations
		annotations := map[string]string{}
		for key, value := range service.Annotations {
			annotations[key] = value
		}
		sc.ObjectMeta.Annotations = annotations

		// fillTemplate fills the pod template with the value calculated from config
		fillTemplate := func(template *api.PodTemplateSpec) {
			template.Spec.Containers[0].Env = envs
			template.Spec.Containers[0].Command = cmds
			template.Spec.Containers[0].WorkingDir = service.WorkingDir
			template.Spec.Containers[0].VolumeMounts = volumesMount
			template.Spec.Volumes = volumes
			// Configure the container privileged mode
			if service.Privileged == true {
				template.Spec.Containers[0].SecurityContext = &api.SecurityContext{
					Privileged: &service.Privileged,
				}
			}
			template.Spec.Containers[0].Ports = ports
			template.ObjectMeta.Labels = labels
			// Configure the container restart policy.
			switch service.Restart {
			case "", "always":
				template.Spec.RestartPolicy = api.RestartPolicyAlways
			case "no":
				template.Spec.RestartPolicy = api.RestartPolicyNever
			case "on-failure":
				template.Spec.RestartPolicy = api.RestartPolicyOnFailure
			default:
				logrus.Fatalf("Unknown restart policy %s for service %s", service.Restart, name)
			}
		}

		// fillObjectMeta fills the metadata with the value calculated from config
		fillObjectMeta := func(meta *api.ObjectMeta) {
			meta.Labels = labels
			meta.Annotations = annotations
		}

		// update supported controller
		for _, obj := range objects {
			updateController(obj, fillTemplate, fillObjectMeta)
		}

		// If ports not provided in configuration we will not make service
		if len(ports) == 0 {
			logrus.Warningf("[%s] Service cannot be created because of missing port.", name)
		} else {
			objects = append(objects, sc)
		}
		allobjects = append(allobjects, objects...)
	}
	return allobjects
}

// PrintList will take the data converted and decide on the commandline attributes given
func PrintList(objects []runtime.Object, opt convertOptions) error {
	f := createOutFile(opt.outFile)
	defer f.Close()

	var err error
	var files []string

	// if asked to print to stdout or to put in single file
	// we will create a list
	if opt.toStdout || f != nil {
		list := &api.List{}
		list.Items = objects

		// version each object in the list
		list.Items, err = ConvertToVersion(list.Items)
		if err != nil {
			return err
		}

		// version list itself
		listVersion := unversioned.GroupVersion{Group: "", Version: "v1"}
		convertedList, err := api.Scheme.ConvertToVersion(list, listVersion)
		if err != nil {
			return err
		}
		data, err := marshal(convertedList, opt.generateYaml)
		if err != nil {
			return fmt.Errorf("Error in marshalling the List: %v", err)
		}
		files = append(files, print("", "", data, opt.toStdout, opt.generateYaml, f))
	} else {
		var file string
		// create a separate file for each provider
		for _, v := range objects {
			data, err := marshal(v, opt.generateYaml)
			if err != nil {
				return err
			}
			switch t := v.(type) {
			case *api.ReplicationController:
				file = print(t.Name, strings.ToLower(t.Kind), data, opt.toStdout, opt.generateYaml, f)
			case *extensions.Deployment:
				file = print(t.Name, strings.ToLower(t.Kind), data, opt.toStdout, opt.generateYaml, f)
			case *extensions.DaemonSet:
				file = print(t.Name, strings.ToLower(t.Kind), data, opt.toStdout, opt.generateYaml, f)
			case *deployapi.DeploymentConfig:
				file = print(t.Name, strings.ToLower(t.Kind), data, opt.toStdout, opt.generateYaml, f)
			case *api.Service:
				file = print(t.Name, strings.ToLower(t.Kind), data, opt.toStdout, opt.generateYaml, f)
			}
			files = append(files, file)

		}
	}
	if opt.createChart {
		generateHelm(opt.inputFile, files)
	}
	return nil
}

// marshal object runtime.Object and return byte array
func marshal(obj runtime.Object, yamlFormat bool) (data []byte, err error) {
	// convert data to yaml or json
	if yamlFormat {
		data, err = yaml.Marshal(obj)
	} else {
		data, err = json.MarshalIndent(obj, "", "  ")
	}
	if err != nil {
		data = nil
	}
	return
}

// Convert all objects in objs to versioned objects
func ConvertToVersion(objs []runtime.Object) ([]runtime.Object, error) {
	ret := []runtime.Object{}

	for _, obj := range objs {

		objectVersion := obj.GetObjectKind().GroupVersionKind()
		version := unversioned.GroupVersion{Group: objectVersion.Group, Version: objectVersion.Version}
		convertedObject, err := api.Scheme.ConvertToVersion(obj, version)
		if err != nil {
			return nil, err
		}
		ret = append(ret, convertedObject)
	}

	return ret, nil
}

func validateFlags(opt convertOptions, singleOutput bool, dabFile, inputFile string) {
	if len(opt.outFile) != 0 && opt.toStdout {
		logrus.Fatalf("Error: --out and --stdout can't be set at the same time")
	}
	if opt.createChart && opt.toStdout {
		logrus.Fatalf("Error: chart cannot be generated when --stdout is specified")
	}
	if opt.replicas < 0 {
		logrus.Fatalf("Error: --replicas cannot be negative")
	}
	if singleOutput {
		count := 0
		if opt.createD {
			count++
		}
		if opt.createDS {
			count++
		}
		if opt.createRC {
			count++
		}
		if opt.createDeploymentConfig {
			count++
		}
		if count > 1 {
			logrus.Fatalf("Error: only one type of Kubernetes controller can be generated when --out or --stdout is specified")
		}
	}
	if len(dabFile) > 0 && len(inputFile) > 0 && inputFile != DefaultComposeFile {
		logrus.Fatalf("Error: compose file and dab file cannot be specified at the same time")
	}
}

// Convert tranforms docker compose or dab file to k8s objects
func Convert(c *cli.Context) {
	inputFile := c.String("file")
	dabFile := c.String("bundle")
	outFile := c.String("out")
	generateYaml := c.BoolT("yaml")
	toStdout := c.BoolT("stdout")
	createD := c.BoolT("deployment")
	createDS := c.BoolT("daemonset")
	createRC := c.BoolT("replicationcontroller")
	createChart := c.BoolT("chart")
	replicas := c.Int("replicas")
	singleOutput := len(outFile) != 0 || toStdout
	createDeploymentConfig := c.BoolT("deploymentconfig")

	// Create Deployment by default if no controller has be set
	if !createD && !createDS && !createRC && !createDeploymentConfig {
		createD = true
	}

	file := inputFile
	if len(dabFile) > 0 {
		file = dabFile
	}

	opt := convertOptions{
		toStdout:               toStdout,
		createD:                createD,
		createRC:               createRC,
		createDS:               createDS,
		createDeploymentConfig: createDeploymentConfig,
		createChart:            createChart,
		generateYaml:           generateYaml,
		replicas:               replicas,
		inputFile:              file,
		outFile:                outFile,
	}

	validateFlags(opt, singleOutput, dabFile, inputFile)

	komposeObject := KomposeObject{}

	if len(dabFile) > 0 {
		komposeObject = loadBundlesFile(dabFile)
	} else {
		komposeObject = loadComposeFile(inputFile)
	}

	// Convert komposeObject to K8S controllers
	objects := komposeConvert(komposeObject, opt)

	// print output to places as needed
	PrintList(objects, opt)
}

func checkUnsupportedKey(service interface{}) {
	s := structs.New(service)
	for _, f := range s.Fields() {
		if f.IsExported() && !f.IsZero() && f.Name() != "Networks" {
			if count, ok := unsupportedKey[f.Name()]; ok && count == 0 {
				logrus.Warningf("Unsupported key %s - ignoring", composeOptions[f.Name()])
				unsupportedKey[f.Name()]++
			}
		}
	}
}

// Either print to stdout or to file/s
func print(name, trailing string, data []byte, toStdout, generateYaml bool, f *os.File) string {

	file := ""
	if generateYaml {
		file = fmt.Sprintf("%s-%s.yaml", name, trailing)
	} else {
		file = fmt.Sprintf("%s-%s.json", name, trailing)
	}
	if toStdout {
		fmt.Fprintf(os.Stdout, "%s\n", string(data))
		return ""
	} else if f != nil {
		// Write all content to a single file f
		if _, err := f.WriteString(fmt.Sprintf("%s\n", string(data))); err != nil {
			logrus.Fatalf("Failed to write %s to file: %v", trailing, err)
		}
		f.Sync()
	} else {
		// Write content separately to each file
		if err := ioutil.WriteFile(file, []byte(data), 0644); err != nil {
			logrus.Fatalf("Failed to write %s: %v", trailing, err)
		}
		logrus.Printf("file %q created", file)
	}
	return file
}

// Up brings up deployment, svc.
func Up(c *cli.Context) {
	fmt.Println("We are going to create Kubernetes deployment and service for your dockerized application. \n" +
		"If you need more kind of controllers, use 'kompose convert' and 'kubectl create -f' instead. \n")

	factory := cmdutil.NewFactory(nil)
	clientConfig, err := factory.ClientConfig()
	if err != nil {
		logrus.Fatalf("Failed to access the Kubernetes cluster. Make sure you have a Kubernetes running: %v", err)
	}
	client := client.NewOrDie(clientConfig)

	inputFile := c.String("file")
	dabFile := c.String("bundle")

	komposeObject := KomposeObject{}
	opt := convertOptions{
		replicas: 1,
		createD:  true,
	}

	validateFlags(opt, false, dabFile, inputFile)

	if len(dabFile) > 0 {
		komposeObject = loadBundlesFile(dabFile)
	} else {
		komposeObject = loadComposeFile(inputFile)
	}

	//Convert komposeObject to K8S controllers
	objects := komposeConvert(komposeObject, opt)
	objects = sortServicesFirst(objects)

	for _, v := range objects {
		switch t := v.(type) {
		case *extensions.Deployment:
			_, err := client.Deployments(api.NamespaceDefault).Create(t)
			if err != nil {
				logrus.Fatalf("Error: '%v' while creating deployment: %s", err, t.Name)
			}
			logrus.Infof("Successfully created deployment: %s", t.Name)
		case *api.Service:
			_, err := client.Services(api.NamespaceDefault).Create(t)
			if err != nil {
				logrus.Fatalf("Error: '%v' while creating service: %s", err, t.Name)
			}
			logrus.Infof("Successfully created service: %s", t.Name)
		}
	}
	fmt.Println("\nApplication has been deployed to Kubernetes. You can run 'kubectl get deployment,svc' for details.")
}

// the objects that we get can be in any order this keeps services first
// according to best practice kubernetes services should be created first
// http://kubernetes.io/docs/user-guide/config-best-practices/
func sortServicesFirst(objs []runtime.Object) []runtime.Object {
	var svc []runtime.Object
	var others []runtime.Object
	var ret []runtime.Object

	for _, obj := range objs {
		if obj.GetObjectKind().GroupVersionKind().Kind == "Service" {
			svc = append(svc, obj)
		} else {
			others = append(others, obj)
		}
	}
	ret = append(ret, svc...)
	ret = append(ret, others...)
	return ret
}

// updateController updates the given object with the given pod template update function and ObjectMeta update function
func updateController(obj runtime.Object, updateTemplate func(*api.PodTemplateSpec), updateMeta func(meta *api.ObjectMeta)) {
	switch t := obj.(type) {
	case *api.ReplicationController:
		if t.Spec.Template == nil {
			t.Spec.Template = &api.PodTemplateSpec{}
		}
		updateTemplate(t.Spec.Template)
		updateMeta(&t.ObjectMeta)
	case *extensions.Deployment:
		updateTemplate(&t.Spec.Template)
		updateMeta(&t.ObjectMeta)
	case *extensions.ReplicaSet:
		updateTemplate(&t.Spec.Template)
		updateMeta(&t.ObjectMeta)
	case *extensions.DaemonSet:
		updateTemplate(&t.Spec.Template)
		updateMeta(&t.ObjectMeta)
	case *deployapi.DeploymentConfig:
		updateTemplate(t.Spec.Template)
		updateMeta(&t.ObjectMeta)
	}
}
