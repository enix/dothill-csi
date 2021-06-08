package main

import (
	"flag"
	"fmt"

	"github.com/enix/dothill-csi/pkg/common"
	"github.com/enix/dothill-csi/pkg/node"
	"k8s.io/klog"
)

var kubeletPath = flag.String("kubeletpath", "", "Kubelet path (deprecated)")
var bind = flag.String("bind", fmt.Sprintf("unix:///var/run/%s/csi.sock", common.PluginName), "RPC bind URI (can be a UNIX socket path or any URI)")

func main() {
	klog.InitFlags(nil)
	flag.Set("logtostderr", "true")
	flag.Parse()
	klog.Infof("starting dothill storage node plugin %s", common.Version)
	node.New(*kubeletPath).Start(*bind)
}
