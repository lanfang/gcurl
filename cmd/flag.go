package cmd

import (
	"github.com/lanfang/gcurl/config"
	"github.com/spf13/cobra"
)

func init() {
	cobra.OnInitialize(func() {
		initConfig()
	})
	rootCmd.PersistentFlags().StringVarP(&config.G_Conf.ProtoFile, "proto", "b", "", "local proto file")
	rootCmd.Flags().StringVarP(&config.G_Conf.Data, "data", "d", "", "request data")
}

func initConfig() {
	//do something
}
