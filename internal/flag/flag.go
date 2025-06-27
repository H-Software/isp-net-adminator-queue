package flag

import (
	"strings"
	"time"
	"fmt"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	// for testing
	FlagHelp bool
	FlagPhp  bool

	// for asyncq
	FlagRedisAddr        string
	FlagGroupGracePeriod time.Duration
	FlagGroupMaxDelay    time.Duration
	FlagGroupMaxSize     int

	ConfigGatewayWifiFqdn string
	ConfigGateway3Fqdn string
)

func init() {

	flag.String("redis-addr", "localhost:6379", "Redis server address")
	flag.Duration("asynq-grace-period", 10*time.Second, "Group grace period")
	flag.Duration("asynq-max-delay", 30*time.Second, "Group max delay")
	flag.Int("asynq-max-size", 3, "Group max size")

	flag.Bool("help", false, "print usage")
	flag.Bool("php-version", false, "print version of PHP")

	flag.CommandLine.MarkHidden("help")

	flag.Parse()
	viper.BindPFlags(flag.CommandLine)

	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)

	// environments variables support
	viper.BindEnv("redis-addr")

	FlagHelp = viper.GetBool("help")
	FlagPhp = viper.GetBool("php-version")

	FlagRedisAddr = viper.GetString("redis-addr")
	FlagGroupGracePeriod = viper.GetDuration("asynq-grace-period")
	FlagGroupMaxDelay = viper.GetDuration("asynq-max-delay")
	FlagGroupMaxSize = viper.GetInt("asynq-max-size")

	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml") // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("config")   // path to look for the config file in

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	ConfigGatewayWifiFqdn = viper.GetString("gateways.gateway-wifi.fqdn")
	ConfigGateway3Fqdn = viper.GetString("gateways.gateway-3.fqdn")

}

func PrintDefaults() {
	flag.PrintDefaults()
}
