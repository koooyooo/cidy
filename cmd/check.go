/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/koooyooo/cidy/internal/fileutil"

	"github.com/spf13/cobra"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		filePath, _ := cmd.Flags().GetString("file")
		jsonOut, _ := cmd.Flags().GetBool("json")

		if filePath != "" {
			// ファイル一括判定
			if len(args) < 1 {
				fmt.Println("Usage: cidy check -f <file> <CIDR>")
				return
			}
			cidrStr := args[0]
			results, err := fileutil.CheckIPsInFile(filePath, cidrStr)
			if err != nil {
				fmt.Println(err)
				return
			}
			if jsonOut {
				b, _ := json.Marshal(results)
				fmt.Println(string(b))
			} else {
				for _, r := range results {
					if !r.Valid {
						fmt.Printf("%s: Invalid IP address\n", r.IP)
					} else if r.Match {
						fmt.Printf("%s: true\n", r.IP)
					} else {
						fmt.Printf("%s: false\n", r.IP)
					}
				}
			}
			return
		}

		// 通常/JSON判定
		if len(args) < 2 {
			fmt.Println("Usage: cidy check <IP> <CIDR>")
			return
		}
		ipStr := args[0]
		cidrStr := args[1]

		ip := net.ParseIP(ipStr)
		if ip == nil {
			fmt.Println("Invalid IP address")
			return
		}
		_, ipnet, err := net.ParseCIDR(cidrStr)
		if err != nil {
			fmt.Println("Invalid CIDR")
			return
		}
		result := ipnet.Contains(ip)
		if jsonOut {
			type out struct {
				IP    string `json:"ip"`
				CIDR  string `json:"cidr"`
				Match bool   `json:"match"`
			}
			o := out{IP: ipStr, CIDR: cidrStr, Match: result}
			b, _ := json.Marshal(o)
			fmt.Println(string(b))
		} else {
			if result {
				fmt.Println("true")
			} else {
				fmt.Println("false")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)

	// -f/--file フラグ
	checkCmd.Flags().StringP("file", "f", "", "IP address list file")
	// --json フラグ
	checkCmd.Flags().BoolP("json", "j", false, "Output result as JSON")
}
