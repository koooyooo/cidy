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

// CheckIPsInFile checks each IP in the file against the given CIDR and prints the result.
func CheckIPsInFile(filePath string, cidrStr string) error {
	_, ipnet, err := net.ParseCIDR(cidrStr)
	if err != nil {
		return fmt.Errorf("Invalid CIDR")
	}
	ips, err := ReadIPList(filePath)
	if err != nil {
		return fmt.Errorf("File read error: %v", err)
	}
	for _, ipStr := range ips {
		ip := net.ParseIP(ipStr)
		if ip == nil {
			fmt.Printf("%s: Invalid IP address\n", ipStr)
			continue
		}
		if ipnet.Contains(ip) {
			fmt.Printf("%s: true\n", ipStr)
		} else {
			fmt.Printf("%s: false\n", ipStr)
		}
	}
	return nil
}
