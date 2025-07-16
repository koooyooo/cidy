package fileutil

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

// readIPList reads IP addresses from a file, one per line.
func ReadIPList(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var ips []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			ips = append(ips, line)
		}
	}
	return ips, scanner.Err()
}

type CheckResult struct {
	IP    string `json:"ip"`
	CIDR  string `json:"cidr"`
	Match bool   `json:"match"`
	Valid bool   `json:"valid"`
}

func CheckIPsInFile(filePath string, cidrStr string) ([]CheckResult, error) {
	_, ipnet, err := net.ParseCIDR(cidrStr)
	if err != nil {
		return nil, fmt.Errorf("Invalid CIDR")
	}
	ips, err := ReadIPList(filePath)
	if err != nil {
		return nil, fmt.Errorf("File read error: %v", err)
	}
	var results []CheckResult
	for _, ipStr := range ips {
		ip := net.ParseIP(ipStr)
		if ip == nil {
			results = append(results, CheckResult{
				IP:    ipStr,
				CIDR:  cidrStr,
				Match: false,
				Valid: false,
			})
			continue
		}
		results = append(results, CheckResult{
			IP:    ipStr,
			CIDR:  cidrStr,
			Match: ipnet.Contains(ip),
			Valid: true,
		})
	}
	return results, nil
}
