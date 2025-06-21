package flag

import (
	"flag"
	"time"
)

var (
	FlagRedisAddr        = flag.String("redis-addr", "localhost:16379", "Redis server address")
	FlagGroupGracePeriod = flag.Duration("grace-period", 10*time.Second, "Group grace period")
	FlagGroupMaxDelay    = flag.Duration("max-delay", 30*time.Second, "Group max delay")
	FlagGroupMaxSize     = flag.Int("max-size", 3, "Group max size")
)

func init() {
	flag.Parse()
}
