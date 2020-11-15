package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/rajatjindal/autocert-reverseproxy/pkg/proxy"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// proxyCmd represents the proxy command
var proxyCmd = &cobra.Command{
	Use:   "proxy",
	Short: "starts the reverse proxy wrapped with autocert",
	Run: func(cmd *cobra.Command, args []string) {
		p, err := proxy.New()
		if err != nil {
			logrus.Fatal(err)
		}

		p.Start()

		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGTERM)
		<-signals
		logrus.Info("Received SIGTERM. Terminating...")
	},
}

func init() {
	rootCmd.AddCommand(proxyCmd)
}
