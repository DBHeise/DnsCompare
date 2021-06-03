package main

import (
	"DnsCompare"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "verify the servers",
	Run: func(cmd *cobra.Command, args []string) {
		verify()
	},
}

func init() {
	rootCmd.AddCommand(verifyCmd)
	verifyCmd.Flags().StringVarP(&targetDomain, "target", "t", "example.com", "Target Domain name to resolve")

	verifyCmd.Flags().StringVarP(&serverDef, "ServerDef", "s", "", "server configuration definitions (i.e. servers.json file)")
	viper.BindPFlag("ServerDef", verifyCmd.Flags().Lookup("ServerDef"))
	verifyCmd.Flags().StringVarP(&rType, "RecordType", "r", "A", "Type of Query to ask for (available options: A, AAAA, NS, CNAME, PTR, TXT, SRV, SOA, SIG)")
	viper.BindPFlag("RecordType", verifyCmd.Flags().Lookup("RecordType"))

}

func handleResults(method string, serverName string, ans []string, err error) {
	if err != nil {
		log.Error().Err(err).Msg("Error")
	} else if ans == nil {
		log.Warn().Str("Method", method).Str("Server", serverName).Interface("Answers", ans).Msg("Not applicable")
	} else {
		log.Debug().Str("Method", method).Str("Server", serverName).Interface("Answers", ans).Msg("Response")
	}
}

func verify() {
	var err error

	servers, err := DnsCompare.ReadServers(viper.GetString("ServerDef"))
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to read server list json file")
	}

	viper.Set("PickMethod", "all")
	DnsCompare.SetDNSType(viper.GetString("RecordType"))

	var resp []string

	for _, server := range servers {
		resp, err = server.ResolveUDP(targetDomain + ".")
		handleResults("udp", server.Name, resp, err)
		resp, err = server.ResolveTCP(targetDomain + ".")
		handleResults("tcp", server.Name, resp, err)
		resp, err = server.ResolveDOT(targetDomain + ".")
		handleResults("dot", server.Name, resp, err)
		resp, err = server.ResolveDOH(targetDomain + ".")
		handleResults("doh", server.Name, resp, err)
	}
}
