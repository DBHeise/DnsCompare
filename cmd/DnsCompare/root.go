package main

import (
	"DnsCompare"
	"fmt"
	"os"
	"sort"
	"sync"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile      string
	targetDomain string
	parallel     bool
	responseMap  map[string][]string
)

var compareCmd = &cobra.Command{
	Use:   "Compare",
	Short: "Compare DNS results from different servers",
	Run: func(cmd *cobra.Command, args []string) {
		compare()
	},
}

func init() {
	responseMap = make(map[string][]string)
	cobra.OnInitialize(initConfig)
	compareCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/config.json)")

	compareCmd.Flags().StringVarP(&targetDomain, "target", "t", "", "Target Domain name to resolve")
	compareCmd.MarkFlagRequired("target")
	viper.BindPFlag("target", compareCmd.Flags().Lookup("target"))

	compareCmd.Flags().BoolVarP(&parallel, "parallel", "p", false, "Query servers in parallel")

}

func initConfig() {
	DnsCompare.InitializeConfig(cfgFile)
	log = DnsCompare.InitializeLogger("DNSCompare")
}

func Execute() {
	if err := compareCmd.Execute(); err != nil {
		log.Error().Err(err).Msg("An Unhandled Error Occured")
		os.Exit(-1)
	}
}

func resolveSingle(target string, server DnsCompare.DNSServer, wg *sync.WaitGroup) []string {
	resp, err := server.ResolveUDP(targetDomain + ".")
	if err != nil {
		log.Error().Err(err).Msg("Error")
	} else {
		if resp != nil {
			responseMap[server.Name] = resp
			log.Info().Str("Server", server.Name).Interface("resp", resp).Msg("Response")
			wg.Done()
			return resp
		}
	}
	wg.Done()
	return nil
}

func compare() {
	serverListFile := viper.GetString("ServerList")
	servers, err := DnsCompare.ReadServers(serverListFile)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to read server list json file")
	}

	var wg sync.WaitGroup
	for _, server := range servers {
		wg.Add(1)
		if parallel {
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
