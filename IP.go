package utils

import (
	"net"
	"regexp"
)

func GetIPAddressesWithHint(hintRegex string) ([]string, error) {
	var ipAddresses []string

	hint := regexp.MustCompile(hintRegex)

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ip := ipnet.IP.To4(); ip != nil {
				if hint.MatchString(ip.String()) {
					ipAddresses = append(ipAddresses, ip.String())
				}
			}
		}
	}

	return ipAddresses, nil
}
