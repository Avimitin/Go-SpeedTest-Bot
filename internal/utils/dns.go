package utils

import "net"

func GetIP(url string) ([]string, error) {
	ipRec, err := net.LookupIP(url)
	if err != nil {
		return nil, err
	}
	var ret []string
	for _, ip := range ipRec {
		if ip != nil {
			ret = append(ret, ip.String())
		}
	}
	return ret, nil
}
