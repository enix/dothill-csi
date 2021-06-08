package main

import (
	"testing"

	"github.com/enix/dothill-csi/pkg/controller"
	"github.com/enix/dothill-csi/pkg/node"
	"github.com/kubernetes-csi/csi-test/pkg/sanity"
)

// Test starts the drivers in background and runs k8s sanity checks
func Test(t *testing.T) {
	controllerSocketPath := "unix:///tmp/controller.sock"
	nodeSocketPath := "unix:///tmp/node.sock"

	ctrl := controller.New()
	node := node.New()

	go ctrl.Start(controllerSocketPath)
	defer ctrl.Stop()

	go node.Start(nodeSocketPath)
	defer node.Stop()

	sanity.Test(t, &sanity.Config{
		Address:                  nodeSocketPath,
		ControllerAddress:        controllerSocketPath,
		SecretsFile:              "./secrets.yml",
		TestVolumeParametersFile: "./config.yml",
	})
}
