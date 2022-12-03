package cmd

import (
	"flag"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	temporal "go.temporal.io/sdk/client"
	"log"
	"ng-receiver/pkg/server"
	"os"
	"strings"
)

var cfgFile string

type config struct {
	ServerConfig serverConfig   `mapstructure:"http"`
	Temporal     temporalConfig `mapstructure:"temporal"`
}

type temporalConfig struct {
	Address   string `mapstructure:"address"`
	Namespace string `mapstructure:"namespace"`
}

type serverConfig struct {
	Port string `mapstructure:"port"`
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "receiver-service",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.receiver-service.yaml)")

	rootCmd.AddCommand(serveHttpCmd)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".receiver-service")
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

var serveHttpCmd = &cobra.Command{
	Use:     "serve-http",
	Short:   "Run the http server listening for messages",
	Aliases: []string{"serve-http"},
	Run: func(cmd *cobra.Command, args []string) {

		flag.Parse()

		cfg := &config{}
		if err := viper.Unmarshal(cfg); err != nil {
			log.Fatal(err)
		}

		if cfg.ServerConfig.Port == "" {
			log.Fatal("Port must be configured")
		}

		//ctx := context.Background()

		//nc, err := createNatsClient(cfg)
		//if err != nil {
		//	fmt.Printf("unable to create connection %s\n", err)
		//	fmt.Printf("nats config: %v\n", cfg.Nats)
		//	return
		//}
		//defer nc.Close()

		//temporalClient, err := setupTemporalClient(cfg)
		//if err != nil {
		//	glog.Warning("Temporal client could not be created: " + err.Error())
		//}
		//defer temporalClient.Close()

		srv := &server.Server{
			Port: cfg.ServerConfig.Port,
			//TemporalClient: &temporalClient,
		}

		err := srv.RegisterHandlersAndServe()
		if err != nil {
			log.Fatal("Could not start http server", err)
		}
	},
}

func setupTemporalClient(cfg *config) (temporal.Client, error) {
	temporalOptions := temporal.Options{
		HostPort:  cfg.Temporal.Address,
		Namespace: cfg.Temporal.Namespace,
	}

	return temporal.Dial(temporalOptions)
}
