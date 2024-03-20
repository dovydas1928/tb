package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"tb"
	helper "github.com/dovydas1928/tb/cmd"
	proto "github.com/dovydas1928/tb/pkg/proto"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "tb-server",
	Short: "Raspberry PI IO server",
	Long:  `A gRPC server that allows you to do IO operations on the Raspberry PI`,
	Run: func(_ *cobra.Command, args []string) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		log.Println("starting tb server")
		srv, lis, err := tb.NewGrpcServerInsecure(viper.GetString("server.host"), viper.GetString("server.port"))
		if err != nil {
			helper.ExitOnError("unable to create server", err)
		}

		common := tb.NewCommonLocal()
		proto.RegisterCommonServer(srv, tb.NewCommonServer(common))

		if viper.GetBool("gpio.enabled") {
			log.Println("adding gpio service")
			gpio := tb.NewGpioLocal()
			if viper.GetBool("gpio.open") {
				log.Println("opening gpio interface")
				err := gpio.Open(ctx)
				helper.ExitOnError("unable to open gpio", err)
			}
			defer gpio.Close(ctx)
			proto.RegisterGpioServer(srv, tb.NewGpioServer(gpio))
		}

		modprobe := viper.GetStringSlice("modprobe")
		if len(modprobe) > 0 {
			log.Printf("running modprobe for %s\n", modprobe)
			for _, mod := range modprobe {
				err := common.Modprobe(ctx, mod)
				if err != nil {
					helper.ExitOnError(fmt.Sprintf("unable modprobe module '%s'", mod), err)
				}
			}
		}

		log.Fatal(srv.Serve(lis))
	},
}

// Execute the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().BoolP("gpio", "g", false, "gpio service enabled")
	rootCmd.PersistentFlags().Bool("gpio_open", false, "open gpio service on start")
	viper.BindPFlag("gpio.enabled", rootCmd.PersistentFlags().Lookup("gpio"))
	viper.BindPFlag("gpio.open", rootCmd.PersistentFlags().Lookup("gpio_open"))

	rootCmd.PersistentFlags().StringSlice("modprobe", []string{}, "modprobe on start")
	viper.BindPFlag("modprobe", rootCmd.PersistentFlags().Lookup("modprobe"))

	rootCmd.PersistentFlags().StringP("host", "s", "0.0.0.0", "server ip")
	rootCmd.PersistentFlags().IntP("port", "p", 8000, "server port")
	viper.BindPFlag("server.host", rootCmd.PersistentFlags().Lookup("host"))
	viper.BindPFlag("server.port", rootCmd.PersistentFlags().Lookup("port"))

	helper.AddConfigCommand(rootCmd)
}

var configFileName = "config"
var configPath = "/etc/tb-server"

func initConfig() {
	viper.AddConfigPath(configPath)
	viper.SetConfigName(configFileName)

	readErr := viper.ReadInConfig()
	if readErr != nil {
		fmt.Fprintf(os.Stderr, "unable to read config file: %v\n", readErr)
	}
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("tb")
}
