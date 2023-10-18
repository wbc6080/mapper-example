package main

import (
	"errors"
	"os"

	"k8s.io/klog/v2"

	"github.com/kubeedge/virtualdevice/device"
	"github.com/wbc6080/mapper-go/pkg/common"
	"github.com/wbc6080/mapper-go/pkg/config"
	"github.com/wbc6080/mapper-go/pkg/grpcserver"
	"github.com/wbc6080/mapper-go/pkg/httpserver"
	"github.com/wbc6080/mapper-go/pkg/util/grpcclient"
	"github.com/wbc6080/mapper-go/pkg/util/parse"
)

func main() {
	var err error
	var c config.Config

	klog.InitFlags(nil)
	defer klog.Flush()

	if err = c.Parse(); err != nil {
		klog.Fatal(err)
		os.Exit(1)
	}
	klog.Infof("config: %+v", c)

	grpcclient.Init(&c)

	// start grpc server
	grpcServer := grpcserver.NewServer(
		grpcserver.Config{
			SockPath: c.GrpcServer.SocketPath,
			Protocol: common.ProtocolCustomized,
		},
		device.NewDevPanel(),
	)

	panel := device.NewDevPanel()
	err = panel.DevInit(&c)
	if err != nil && !errors.Is(err, parse.ErrEmptyData) {
		klog.Fatal(err)
	}
	klog.Infoln("devInit finished")

	// register to edgecore
	// if dev init mode is register, mapper's dev will init when registry to edgecore
	if c.DevInit.Mode != common.DevInitModeRegister {
		klog.Infoln("======dev init mode is not register, will register to edgecore")
		if _, _, err = grpcclient.RegisterMapper(&c, false); err != nil {
			klog.Fatal(err)
		}
		klog.Infoln("registerMapper finished")
	}
	go panel.DevStart()

	httpServer := httpserver.NewRestServer(panel)
	go httpServer.StartServer()

	defer grpcServer.Stop()
	if err = grpcServer.Start(); err != nil {
		klog.Fatal(err)
	}
}
