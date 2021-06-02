package DnsCompare

import (
	"fmt"
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

func init() {

	viper.SetDefault("Log.ReportCaller", false)
	viper.SetDefault("Log.Output", "console")
	viper.SetDefault("Log.Level", "trace")
}

// InitializeConfig - Setup viper & config
func InitializeConfig(cfgFile string) {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {

		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
		viper.AddConfigPath(home)

		currentPath, err := filepath.Abs(".")
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
		viper.AddConfigPath(currentPath)

		parentPath, err := filepath.Abs("..")
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
		viper.AddConfigPath(parentPath)

		viper.SetConfigName("config.json")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	viper.ReadInConfig()
	//if err := viper.ReadInConfig(); err == nil {
	//fmt.Println("Using config file:", viper.ConfigFileUsed())
	//} else {
	//panic(err)
	//}

}
func GenerateDefaultConfig(outfile string) error {
	return viper.WriteConfigAs(outfile)
}
