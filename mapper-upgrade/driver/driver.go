package driver

import (
	"fmt"
	"k8s.io/klog/v2"
	"math/rand"
	"sync"
)

func NewClient(protocol ProtocolConfig) (*CustomizedClient, error) {
	client := &CustomizedClient{
		ProtocolConfig: protocol,
		deviceMutex:    sync.Mutex{},
		// TODO initialize the variables you added
	}
	return client, nil
}

func (c *CustomizedClient) InitDevice() error {
	// TODO: add init operation
	// you can use c.ProtocolConfig and c.ProtocolCommonConfig
	klog.Infof("Init device%d successful, protocolID: %v", c.DeviceID, c.ProtocolID)
	klog.Infof("I can get Info: %v %v ", c.SerialPort, c.BaudRate)
	return nil
}

func (c *CustomizedClient) GetDeviceData(visitor *VisitorConfig) (interface{}, error) {
	// TODO: get device's data
	// you can use c.ProtocolConfig,c.ProtocolCommonConfig and visitor
	if visitor.VisitorConfigData.DataType == "int" {
		if c.intMaxValue <= 0 {
			return nil, fmt.Errorf("max value is %d, should > 0", c.intMaxValue)
		}
		return rand.Intn(c.intMaxValue), nil
	} else if visitor.DataType == "float" {
		return rand.Float64(), nil
	} else {
		return nil, fmt.Errorf("unrecognized data type: %s", visitor.DataType)
	}
	return nil, nil
}

func (c *CustomizedClient) SetDeviceData(data interface{}, visitor *VisitorConfig) error {
	// TODO: set device's data
	// you can use c.ProtocolConfig,c.ProtocolCommonConfig and visitor
	if visitor.DataType == "int" {
		c.intMaxValue = int(data.(int64))
	} else {
		return fmt.Errorf("unrecognized data type: %s", visitor.DataType)
	}
	return nil
}

func (c *CustomizedClient) StopDevice() error {
	// TODO: stop device
	// you can use c.ProtocolConfig and c.ProtocolCommonConfig
	klog.Infof("Stop device%d successful", c.DeviceID)
	return nil
}
