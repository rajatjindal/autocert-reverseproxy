package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/rajatjindal/autocert-reverseproxy/pkg/proxy"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var proxyServer = proxy.AutocertServer{}

// proxyCmd represents the proxy command
var proxyCmd = &cobra.Command{
	Use:   "proxy",
	Short: "starts the reverse proxy wrapped with autocert",
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

func init() {
	rootCmd.AddCommand(proxyCmd)

	proxyCmd.Flags().StringVar(&proxyServer.CertCacheDir, "cert-cache-dir", "", "dir to cache certs to")
	proxyCmd.Flags().StringArrayVar(&proxyServer.AllowedHost, "allowed-host", []string{}, "hostnames allowed to request cert")
	proxyCmd.Flags().StringVar(&proxyServer.HTTPAddr, "http-addr", ":80", "http address")
	proxyCmd.Flags().StringVar(&proxyServer.HTTPSAddr, "https-addr", ":443", "https address")
	proxyCmd.Flags().StringVar(&proxyServer.Upstream, "upstream", "http://localhost:8080", "upstream service")
}
