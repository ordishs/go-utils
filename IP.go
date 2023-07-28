package utils

import (
	"context"
	"net"
	"regexp"
	"time"
)

func GetIPAddressesWithHint(hintRegex string) ([]string, error) {
	var ipAddresses []string

	hint, err := regexp.Compile(hintRegex)
	if err != nil {
		return nil, err
	}

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

func GetPublicIPAddress() (string, error) {
	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Millisecond * time.Duration(10000),
			}
			return d.DialContext(ctx, network, "resolver1.opendns.com:53")
		},
	}
	ip, err := r.LookupHost(context.Background(), "myip.opendns.com")
	if err != nil {
		return "", err
	}

	return ip[0], nil
}
