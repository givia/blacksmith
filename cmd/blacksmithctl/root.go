package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	httptransport "github.com/go-openapi/runtime/client"

	"github.com/cafebazaar/blacksmith/swagger/client"
	"github.com/go-openapi/strfmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:           "blacksmithctl",
	Short:         "blacksmithctl is a client for blacksmith",
	Long:          "blacksmithctl is a client for blacksmith",
	SilenceUsage:  true,
	SilenceErrors: true,
}

func newSwaggerClient() *client.Salesman {
	var host string
	var ok bool

	if host, ok = viper.Get("host").(string); !ok {
		log.Fatalf("%q is not set in config", "host")
	}

	tlsClient, err := httptransport.TLSClient(httptransport.TLSClientOptions{
		ServerName:  "localhost",
		Certificate: "./certs/cert.pem",
		Key:         "./certs/key.pem",
		CA:          "./certs/cert.pem",
	})
	if err != nil {
		log.Fatal("Error creating TLSClient:", err)
	}

	transport := httptransport.NewWithClient(
		host,
		client.DefaultBasePath,
		client.DefaultSchemes,
		tlsClient,
	)
	return client.New(transport, strfmt.NewFormats())
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	RootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	RootCmd.PersistentFlags().IntP("timeout", "t", 2, "timeout in seconds")
	RootCmd.PersistentFlags().StringP("output", "o", "yaml", "output format. One of: yaml|json")

	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName("config")                // name of config file (without extension)
	viper.AddConfigPath("$HOME/.blacksmithctl/") // adding home directory as first search path
	viper.AutomaticEnv()                         // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}
	if verbose, err := RootCmd.PersistentFlags().GetBool("verbose"); err == nil && verbose {
		fmt.Println("Using config file", viper.ConfigFileUsed())
	}
}

func checkArgs(args, mandatoryArgs []string) error {
	if got, want := len(args), len(mandatoryArgs); got != want {
		if got < want {
			isOrAre := "are"
			if want-got == 1 {
				isOrAre = "is"
			}
			return fmt.Errorf("%s %s missing", strings.Join(mandatoryArgs[got:], " "), isOrAre)
		} else if got > want {
			return fmt.Errorf("%d extra args was given", got-want)
		}
	}
	return nil
}
