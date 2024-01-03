package util

import (
	"fmt"
	"net"
	"os"
	"strings"

	"opendev.org/airship/kubernetes-entrypoint/logger"
	"opendev.org/airship/kubernetes-entrypoint/util/env"
)

func GetIp() (string, error) {
	var iface string
	if iface = os.Getenv("INTERFACE_NAME"); iface == "" {
		return "", fmt.Errorf("environment variable INTERFACE_NAME not set")
	}
	i, err := net.InterfaceByName(iface)
	if err != nil {
		return "", fmt.Errorf("cannot get iface: %v", err)
	}

	address, err := i.Addrs()
	if err != nil || len(address) == 0 {
		return "", fmt.Errorf("cannot get ip: %v", err)
	}
	//Take first element to get rid of subnet
	ip := strings.Split(address[0].String(), "/")[0]
	return ip, nil
}

func ContainsSeparator(envString string, kind string) bool {
	if strings.Contains(envString, env.Separator) {
		logger.Error.Printf("%s doesn't accept namespace: %s", kind, envString)
		return true
	}
	return false
}
