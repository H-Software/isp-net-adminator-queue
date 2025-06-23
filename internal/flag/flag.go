package flag

import (
	"strings"
	"time"

	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	FlagRedisAddr string
	FlagGroupGracePeriod time.Duration
	FlagGroupMaxDelay time.Duration
	FlagGroupMaxSize  int
)

func init() {

	flag.String("redis-addr", "localhost:6379", "Redis server address")
	flag.Duration("grace-period", 10*time.Second, "Group grace period")
	flag.Duration("max-delay", 30*time.Second, "Group max delay")
	flag.Int("max-size", 3, "Group max size")

	flag.Parse()
	viper.BindPFlags(flag.CommandLine)

	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)

	viper.BindEnv("redis-addr")

	FlagRedisAddr = viper.GetString("redis-addr")
	FlagGroupGracePeriod = viper.GetDuration("grace-period")
	FlagGroupMaxDelay = viper.GetDuration("max-delay")
 	FlagGroupMaxSize = viper.GetInt("max-size")
}
