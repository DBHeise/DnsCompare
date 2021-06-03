package main

import (
	"DnsCompare"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	targetDomain string
	serverDef    string
	mode         string
	rType        string
	parallel     bool
	pickMethod   string
	responseMap  map[string][]string
)

var compareCmd = &cobra.Command{
	Use:   "compare",
	Short: "Compare DNS results from different servers",
	Run: func(cmd *cobra.Command, args []string) {
		compare()
	},
}

func init() {
	responseMap = make(map[string][]string)
	rootCmd.AddCommand(compareCmd)

	compareCmd.Flags().StringVarP(&targetDomain, "target", "t", "", "Target Domain name to resolve")
	compareCmd.MarkFlagRequired("target")

	compareCmd.Flags().StringVarP(&serverDef, "ServerDef", "s", "./servers.json", "server configuration definitions (i.e. servers.json file)")
	viper.BindPFlag("ServerDef", compareCmd.Flags().Lookup("ServerDef"))
	compareCmd.Flags().StringVarP(&mode, "DNSMode", "m", "UDP", "DNS mode to use (available options: UDP, TCP, DOT, DOH)")
	viper.BindPFlag("DNSMode", compareCmd.Flags().Lookup("DNSMode"))
	compareCmd.Flags().BoolVarP(&parallel, "Parallel", "p", false, "Query servers in parallel")
	viper.BindPFlag("Parallel", compareCmd.Flags().Lookup("Parallel"))
	compareCmd.Flags().StringVarP(&rType, "RecordType", "r", "A", "Type of Query to ask for (available options: A, AAAA, NS, CNAME, PTR, TXT, SRV, SOA, SIG)")
	viper.BindPFlag("RecordType", compareCmd.Flags().Lookup("RecordType"))
	compareCmd.Flags().StringVarP(&pickMethod, "PickMethod", "k", "first4", "Method to choose a server address (available options: first,random,first4,first6,random4,random6,all,all4,all6)")
	viper.BindPFlag("PickMethod", compareCmd.Flags().Lookup("PickMethod"))
}

func resolveSingle(target string, server DnsCompare.DNSServer, wg *sync.WaitGroup) []string {
	var resp []string
	var err error

	switch strings.ToLower(viper.GetString("DNSMode")) {
	case "udp":
		resp, err = server.ResolveUDP(targetDomain + ".")
	case "tcp":
		resp, err = server.ResolveTCP(targetDomain + ".")
	case "dot":
		resp, err = server.ResolveDOT(targetDomain + ".")
	case "doh":
		resp, err = server.ResolveDOH(targetDomain + ".")
	default:
		log.Fatal().Msg("Invalid mode specified! (Available options: udp,tcp,dot,doh)")
	}

	if err != nil {
		log.Error().Err(err).Msg("Error")
	} else {
		if resp != nil {
			responseMap[server.Name] = resp
			log.Debug().Str("Server", server.Name).Interface("resp", resp).Msg("Response")
			wg.Done()
			return resp
		}
	}
	wg.Done()
	return nil
}

func compare() {
	log.Trace().Msg("Compare")

	servers, err := DnsCompare.ReadServers(viper.GetString("ServerDef"))
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to read server list json file")
	}

	DnsCompare.SetDNSType(viper.GetString("RecordType"))

	log.Trace().Int("ServerCount", len(servers)).Msg("Iterating Servers")
	var wg sync.WaitGroup
	for _, server := range servers {
		wg.Add(1)
		if viper.GetBool("Parallel") {
			go resolveSingle(targetDomain, server, &wg)
		} else {
			resolveSingle(targetDomain, server, &wg)
		}
	}
	wg.Wait()

	//Pretty Print
	longestName := 0
	keys := make([]string, 0, len(responseMap))
	for k := range responseMap {
		keys = append(keys, k)
		keyLength := len(k)
		if keyLength > longestName {
			longestName = keyLength
		}
	}
	sort.Strings(keys)
	for _, name := range keys {
		answers := responseMap[name]
		sort.Strings(answers)
		fmt.Printf("%-*s %s\n", longestName, name, answers)
	}
}
