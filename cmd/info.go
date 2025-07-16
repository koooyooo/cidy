/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/spf13/cobra"
)

var infoJSON bool

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info <CIDR>",
	Short: "Show information about a CIDR range",
	Long: `Show detailed information about a given CIDR range, including network address, broadcast address, subnet mask, usable host range, and host count.

Examples:
  cidy info 192.168.0.0/24
  cidy info --json 192.168.0.0/24
`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cidr := args[0]
		ip, ipnet, err := net.ParseCIDR(cidr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid CIDR: %v\n", err)
			os.Exit(1)
		}

		network := ip.Mask(ipnet.Mask)
		broadcast := getBroadcastAddress(ipnet)
		subnetMask := net.IP(ipnet.Mask).String()
		firstIP, lastIP := getHostRange(ipnet)
		usableHosts := countUsableHosts(ipnet)

		if infoJSON {
			result := map[string]interface{}{
				"cidr":              cidr,
				"network_address":   network.String(),
				"broadcast_address": broadcast.String(),
				"subnet_mask":       subnetMask,
				"host_range": map[string]string{
					"from": firstIP.String(),
					"to":   lastIP.String(),
				},
				"usable_hosts": usableHosts,
			}
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			enc.Encode(result)
		} else {
			fmt.Printf("CIDR:                %s\n", cidr)
			fmt.Printf("Network Address:     %s\n", network.String())
			fmt.Printf("Broadcast Address:   %s\n", broadcast.String())
			fmt.Printf("Subnet Mask:         %s\n", subnetMask)
			fmt.Printf("Host Range:          %s - %s\n", firstIP.String(), lastIP.String())
			fmt.Printf("Usable Hosts:        %d\n", usableHosts)
		}
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
	infoCmd.Flags().BoolVarP(&infoJSON, "json", "j", false, "Output result in JSON format")
}

// getBroadcastAddress returns the broadcast address for the given IPNet
func getBroadcastAddress(ipnet *net.IPNet) net.IP {
	ip := ipnet.IP.To4()
	if ip == nil {
		ip = ipnet.IP.To16()
	}
	mask := ipnet.Mask
	broadcast := make(net.IP, len(ip))
	for i := 0; i < len(ip); i++ {
		broadcast[i] = ip[i] | ^mask[i]
	}
	return broadcast
}

// getHostRange returns the first and last usable IP addresses in the subnet
func getHostRange(ipnet *net.IPNet) (net.IP, net.IP) {
	network := ipnet.IP.Mask(ipnet.Mask)
	broadcast := getBroadcastAddress(ipnet)
	first := make(net.IP, len(network))
	last := make(net.IP, len(broadcast))
	copy(first, network)
	copy(last, broadcast)

	// For IPv4, increment/decrement for usable range
	if first.To4() != nil {
		first[3]++
		last[3]--
	} else {
		// For IPv6, just return network and broadcast (no usable host concept)
	}
	return first, last
}

// countUsableHosts returns the number of usable hosts in the subnet
func countUsableHosts(ipnet *net.IPNet) int {
	ones, bits := ipnet.Mask.Size()
	if bits != 32 {
		// For IPv6, return 0 (no usable host concept)
		return 0
	}
	if ones >= 31 {
		// /31 and /32 have no usable hosts
		return 0
	}
	return 1<<(32-ones) - 2
}
