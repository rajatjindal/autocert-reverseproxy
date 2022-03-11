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

var (
	upstreamFile string
	certCacheDir string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "autocert-reverseproxy",
	Short: "golang reverse proxy that handles cert provisioning as well",
	Run: func(cmd *cobra.Command, args []string) {
		proxyServer, err := proxy.New("certs", "upstream.yaml")
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
	rootCmd.Flags().StringVar(&certCacheDir, "cert-cache-dir", "certs", "dir to cache certs to")
	rootCmd.Flags().StringVar(&upstreamFile, "upstream-file", "upstream.yaml", "upstream map file")
}
