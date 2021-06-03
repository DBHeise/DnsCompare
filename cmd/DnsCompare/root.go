package main

import (
	"DnsCompare"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	loglevel  string
	logcolors bool
	logformat string
	logoutput string
	logtime   string
	logreport bool
)

var rootCmd = &cobra.Command{}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&loglevel, "Log.Level", "info", "Log Level (available options: panic,fatal,error,warn,info,debug,trace)")
	viper.BindPFlag("Log.Level", rootCmd.PersistentFlags().Lookup("Log.Level"))
	rootCmd.PersistentFlags().BoolVar(&logcolors, "Log.Colors", true, "Use colors in logging (only valid for console logging)")
	viper.BindPFlag("Log.Level", rootCmd.PersistentFlags().Lookup("Log.Level"))
	rootCmd.PersistentFlags().StringVar(&logformat, "Log.Format", "info", "Log format (available options: plain, json)")
	viper.BindPFlag("Log.Format", rootCmd.PersistentFlags().Lookup("Log.Format"))
	rootCmd.PersistentFlags().StringVar(&logoutput, "Log.Output", "console", "Log output (available options: console, {filename})")
	viper.BindPFlag("Log.Output", rootCmd.PersistentFlags().Lookup("Log.Output"))
	rootCmd.PersistentFlags().StringVar(&logtime, "Log.TimeFormat", "RFC3339", "Log timestamp format (available options: ansi=ansic,unix=unixdate,ruby=rubydate,rfc822,rfc822z,rfc850,rfc1123,rfc1123z,json=rfc3339,rfc3339nano,kitchen,stamp,stampmili,stampmicro,stampnano,{golang time format string})")
	viper.BindPFlag("Log.TimeFormat", rootCmd.PersistentFlags().Lookup("Log.TimeFormat"))
	rootCmd.PersistentFlags().BoolVar(&logreport, "Log.ReportCaller", false, "Report the calling function to the log")
	viper.BindPFlag("Log.ReportCaller", rootCmd.PersistentFlags().Lookup("Log.ReportCaller"))
}

func initConfig() {
	DnsCompare.InitializeConfig()
	log = DnsCompare.InitializeLogger()
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Error().Err(err).Msg("An Unhandled Error Occured")
		os.Exit(-1)
	}
}
