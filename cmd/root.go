package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/rajatjindal/autocert-reverseproxy/pkg/proxy"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var proxyServer = proxy.AutocertServer{}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "autocert-reverseproxy",
	Short: "golang reverse proxy that handles cert provisioning as well",
	Run: func(cmd *cobra.Command, args []string) {
		err := proxyServer.Init()
		if err != nil {
			logrus.Fatal(err)
		}

		proxyServer.Start()

		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGTERM)
		<-signals
		logrus.Info("Received SIGTERM. Terminating...")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVar(&proxyServer.CertCacheDir, "cert-cache-dir", "", "dir to cache certs to")
	rootCmd.Flags().StringArrayVar(&proxyServer.AllowedHost, "allowed-host", []string{}, "hostnames allowed to request cert")
	rootCmd.Flags().StringVar(&proxyServer.HTTPAddr, "http-addr", ":80", "http address")
	rootCmd.Flags().StringVar(&proxyServer.HTTPSAddr, "https-addr", ":443", "https address")
	rootCmd.Flags().StringVar(&proxyServer.Upstream, "upstream", "http://localhost:8080", "upstream service")
}
