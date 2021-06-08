package DnsCompare

import (
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault("Log.Level", "info") //panic,fatal,error,warn,info.debug,trace
	viper.SetDefault("Log.Colors", true)
	viper.SetDefault("Log.Format", "plain")       //plain, json
	viper.SetDefault("Log.Output", "console")     //console, {file}
	viper.SetDefault("Log.TimeFormat", "RFC3339") //ansi=ansic,unix=unixdate,ruby=rubydate,rfc822,rfc822z,rfc850,rfc1123,rfc1123z,json=rfc3339,rfc3339nano,kitchen,stamp,stampmili,stampmicro,stampnano,{golang time format string}
	viper.SetDefault("Log.ReportCaller", false)

	/*
		viper.SetDefault("ServerDef", "./servers.json")
		viper.SetDefault("DNSMode", "UDP")
		viper.SetDefault("Parallel", false)
		viper.SetDefault("RecordType", "A")
		viper.SetDefault("PickMethod", "first") //first,random,first4,first6,random4,random6,all,all4,all6
	*/
}

// InitializeConfig - Setup viper & config
func InitializeConfig() {
	home, err := homedir.Dir()
	cobra.CheckErr(err)
	viper.AddConfigPath(home)

	currentPath, err := filepath.Abs(".")
	cobra.CheckErr(err)
	viper.AddConfigPath(currentPath)

	parentPath, err := filepath.Abs("..")
	cobra.CheckErr(err)
	viper.AddConfigPath(parentPath)

	viper.SetConfigName("config.json")

	viper.AutomaticEnv()
	viper.ReadInConfig()
}

func GenerateDefaultConfig(outfile string) error {
	return viper.WriteConfigAs(outfile)
}
