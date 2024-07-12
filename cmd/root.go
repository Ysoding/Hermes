package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Ysoding/Hermes/proxy"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "hermes",
	Short: "Hermes is a proxy tool",
	Long:  `Hermes is a proxy tool`,
	Run: func(cmd *cobra.Command, args []string) {
		proxy := &proxy.Proxy{ListenAddr: ":9999"}
		http.ListenAndServe(proxy.ListenAddr, proxy)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
