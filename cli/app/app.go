/*
Copyright 2016 The Kubernetes Authors All rights reserved.

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
	"os/exec"
	"fmt"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/urfave/cli"

	// install kubernetes api
	_ "k8s.io/kubernetes/pkg/api/install"
	_ "k8s.io/kubernetes/pkg/apis/extensions/install"

	// install OpenShift api
	_ "github.com/openshift/origin/pkg/deploy/api/install"
	_ "github.com/openshift/origin/pkg/image/api/install"
	_ "github.com/openshift/origin/pkg/route/api/install"

	"github.com/kubernetes-incubator/kompose/pkg/kobject"
	"github.com/kubernetes-incubator/kompose/pkg/loader"
	"github.com/kubernetes-incubator/kompose/pkg/transformer"
	"github.com/kubernetes-incubator/kompose/pkg/transformer/kubernetes"
	"github.com/kubernetes-incubator/kompose/pkg/transformer/openshift"
)

const (
	DefaultComposeFile = "docker-compose.yml"
	DefaultProvider    = "kubernetes"
)

var inputFormat = "compose"

func validateFlags(c *cli.Context, opt *kobject.ConvertOptions) {

	if opt.OutFile == "-" {
		opt.ToStdout = true
		opt.OutFile = ""
	}

	if len(opt.OutFile) != 0 && opt.ToStdout {
		logrus.Fatalf("Error: --out and --stdout can't be set at the same time")
	}

	if opt.CreateChart && opt.ToStdout {
		logrus.Fatalf("Error: chart cannot be generated when --stdout is specified")
	}

	if opt.Replicas < 0 {
		logrus.Fatalf("Error: --replicas cannot be negative")
	}

	dabFile := c.GlobalString("bundle")

	if len(dabFile) > 0 {
		inputFormat = "bundle"
		opt.InputFile = dabFile
	}

	if len(dabFile) > 0 && c.GlobalIsSet("file") {
		logrus.Fatalf("Error: 'compose' file and 'dab' file cannot be specified at the same time")
	}

	if len(c.Args()) != 0 {
		logrus.Fatal("Unknown Argument(s): ", strings.Join(c.Args(), ","))
	}
}

func validateControllers(opt *kobject.ConvertOptions) {

	singleOutput := len(opt.OutFile) != 0 || opt.OutFile == "-" || opt.ToStdout

	if opt.Provider == "kubernetes" {
		// create deployment by default if no controller has been set
		if !opt.CreateD && !opt.CreateDS && !opt.CreateRC {
			opt.CreateD = true
		}
		if singleOutput {
			count := 0
			if opt.CreateD {
				count++
			}
			if opt.CreateDS {
				count++
			}
			if opt.CreateRC {
				count++
			}
			if count > 1 {
				logrus.Fatalf("Error: only one kind of Kubernetes resource can be generated when --out or --stdout is specified")
			}
		}

	} else if opt.Provider == "openshift" {
		// create deploymentconfig by default if no controller has been set
		if !opt.CreateDeploymentConfig {
			opt.CreateDeploymentConfig = true
		}
		if singleOutput {
			count := 0
			if opt.CreateDeploymentConfig {
				count++
			}
			// Add more controllers here once they are available in OpenShift
			// if opt.foo {count++}

			if count > 1 {
				logrus.Fatalf("Error: only one kind of OpenShift resource can be generated when --out or --stdout is specified")
			}
		}
	}
}

// Convert transforms docker compose or dab file to k8s objects
func Convert(c *cli.Context) {
	opt := kobject.ConvertOptions{
		ToStdout:               c.BoolT("stdout"),
		CreateChart:            c.BoolT("chart"),
		GenerateYaml:           c.BoolT("yaml"),
		Replicas:               c.Int("replicas"),
		InputFile:              c.GlobalString("file"),
		OutFile:                c.String("out"),
		Provider:               strings.ToLower(c.GlobalString("provider")),
		CreateD:                c.BoolT("deployment"),
		CreateDS:               c.BoolT("daemonset"),
		CreateRC:               c.BoolT("replicationcontroller"),
		CreateDeploymentConfig: c.BoolT("deploymentconfig"),
		EmptyVols:              c.BoolT("emptyvols"),
	}

	validateFlags(c, &opt)
	validateControllers(&opt)

	// loader parses input from file into komposeObject.
	l, err := loader.GetLoader(inputFormat)
	if err != nil {
		logrus.Fatal(err)
	}

	komposeObject := kobject.KomposeObject{
		ServiceConfigs: make(map[string]kobject.ServiceConfig),
	}
	komposeObject = l.LoadFile(opt.InputFile)

	// Get a transformer that maps komposeObject to provider's primitives
	t := getTransformer(opt)

	// Do the transformation
	objects := t.Transform(komposeObject, opt)

	// Print output
	kubernetes.PrintList(objects, opt)
}

// Up brings up deployment, svc.
func Up(c *cli.Context) {
	opt := kobject.ConvertOptions{
		InputFile: c.GlobalString("file"),
		Replicas:  1,
		Provider:  strings.ToLower(c.GlobalString("provider")),
		EmptyVols: c.BoolT("emptyvols"),
	}
	validateFlags(c, &opt)
	validateControllers(&opt)

	var hasMinikube bool
	minikubePath, err := exec.LookPath("minikube"); if err != nil {
		hasMinikube = false
	} else {
		hasMinikube = true
	}

	// loader parses input from file into komposeObject.
	l, err := loader.GetLoader(inputFormat)
	if err != nil {
		logrus.Fatal(err)
	}

	komposeObject := kobject.KomposeObject{
		ServiceConfigs: make(map[string]kobject.ServiceConfig),
	}
	komposeObject = l.LoadFile(opt.InputFile)

	// Get the transformer
	t := getTransformer(opt)

	//Submit objects to provider
	errDeploy := t.Deploy(komposeObject, opt)

	if errDeploy != nil && opt.Provider == "kubernetes" && hasMinikube {
		fmt.Printf("Unable to deploy using kubernetes, but found minikube. Should we use minikube to deploy? (y/n): ")
		c := 'N'
		_, _ = fmt.Scanf("%c", &c)

		if c == 'y' || c == 'Y' {
			/* Spawn minikube, and then retry connecting via kubernetes */
			fmt.Printf("Bringing up minikube...")
			minikubeCmd := exec.Command(minikubePath, "start")
			minikubeOutput, err := minikubeCmd.CombinedOutput()
			if err != nil {
				logrus.Fatalf("Error spawning minikube: %s", err)
			}
			fmt.Printf("%s", minikubeOutput)

			/* Need to pause slightly to give apiServer a chance to start */
			time.Sleep(5 * time.Second)

			errDeploy = t.Deploy(komposeObject, opt)
		}
	}

	if errDeploy != nil {
		logrus.Fatalf("Error while deploying application: %s", errDeploy)
	}
}

// Down deletes all deployment, svc.
func Down(c *cli.Context) {
	opt := kobject.ConvertOptions{
		InputFile: c.GlobalString("file"),
		Replicas:  1,
		Provider:  strings.ToLower(c.GlobalString("provider")),
		EmptyVols: c.BoolT("emptyvols"),
	}
	validateFlags(c, &opt)
	validateControllers(&opt)

	// loader parses input from file into komposeObject.
	l, err := loader.GetLoader(inputFormat)
	if err != nil {
		logrus.Fatal(err)
	}

	komposeObject := kobject.KomposeObject{
		ServiceConfigs: make(map[string]kobject.ServiceConfig),
	}
	komposeObject = l.LoadFile(opt.InputFile)

	// Get the transformer
	t := getTransformer(opt)

	//Remove deployed application
	errUndeploy := t.Undeploy(komposeObject, opt)
	if errUndeploy != nil {
		logrus.Fatalf("Error while deleting application: %s", errUndeploy)
	}

}

// Convenience method to return the appropriate Transformer based on
// what provider we are using.
func getTransformer(opt kobject.ConvertOptions) transformer.Transformer {
	var t transformer.Transformer
	if opt.Provider == "kubernetes" {
		// Create/Init new Kubernetes object with CLI opts
		t = &kubernetes.Kubernetes{Opt: opt}
	} else {
		// Create/Init new OpenShift object that is initialized with a newly
		// created Kubernetes object. Openshift inherits from Kubernetes
		t = &openshift.OpenShift{Kubernetes: kubernetes.Kubernetes{Opt: opt}}
	}
	return t
}

func askForConfirmation() bool {
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		logrus.Fatal(err)
	}
	if response == "yes" {
		return true
	} else if response == "no" {
		return false
	} else {
		fmt.Println("Please type yes or no and then press enter:")
		return askForConfirmation()
	}
}
