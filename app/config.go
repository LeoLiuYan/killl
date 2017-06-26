package app

import (
	"github.com/ogier/pflag"
)

const (
	whitelistKey      = "KILLL_WHITELIST"
	labelsKey         = "KILLL_LABELS"
	imageCleanKey     = "KILLL_IMAGE_CLEAN"
	imageBeforeKey    = "KILLL_IMAGE_BEFORE"
	containerCheckKey = "KILLL_CONTAINER_CHECK"
	verboseKey        = "KILLL_VERBOSE"
	dockerHostKey     = "KILLL_DOCKER_HOST"
	dockerAPIKey      = "kILLL_DOCKER_API"
	requestTimeoutKey = "KILLL_REQUEST_TIMEOUT"
	keepAliveKey      = "KILLL_KEEPALIVE"
	ringbufferKey     = "KILLL_RING_BUFFER"
	specificPrefixKey = "KILLL_MESOS_PREFIX"
)

type Config struct {
	whitelist        StringList
	containerCheck   Interval
	imageClean       Interval
	imageBefore      Interval
	labels           StringList
	Verbose          bool
	dockerHost       string
	dockerAPIVersion string
	requestTimeout   Interval
	keepAlive        Interval
	ringBuffer       uint32
	specificPrefix   string
}

func NewConfig() (cfg *Config) {
	cfg = new(Config)

	// whitelist
	cfg.whitelist = envStringList(whitelistKey, "")
	// labels
	cfg.labels = envStringList(labelsKey, "")
	// container check
	cfg.containerCheck = envDuration(containerCheckKey, "10m")
	// image clean
	cfg.imageClean = envDuration(imageCleanKey, "1h")
	// image befoe
	cfg.imageBefore = envDuration(imageBeforeKey, "1440h")
	// verbose
	cfg.Verbose = true
	// docker host
	cfg.dockerHost = "unix:///var/run/docker.sock"
	// docker api version
	cfg.dockerAPIVersion = "v1.28"
	// request timeout
	cfg.requestTimeout = envDuration(requestTimeoutKey, "5s")
	// keepalive
	cfg.keepAlive = envDuration(keepAliveKey, "30s")
	// ring buffer
	cfg.ringBuffer = 512
	// mesos prefix
	cfg.specificPrefix = "mesos"

	return
}

func (cfg *Config) AddFlagSet(flagSet *pflag.FlagSet) {
	flagSet.VarP(&cfg.whitelist, "whitelist", "w", "whitelist for containers not kill by this app, by container name")
	flagSet.VarP(&cfg.labels, "labels", "l", "labels for label the container which should run")
	flagSet.Var(&cfg.imageClean, "image_clean", "interval for clean old image")
	flagSet.Var(&cfg.imageBefore, "image_before", "clean image before this time, unit: h, m, s")
	flagSet.Var(&cfg.containerCheck, "container_check", "interval for check container")
	flagSet.StringVar(&cfg.dockerHost, "docker_host", env(dockerHostKey, "unix:///var/run/docker.sock"), "docker host")
	flagSet.StringVar(&cfg.dockerAPIVersion, "docker_api", env(dockerAPIKey, "v1.28"), "docker api version")
	flagSet.Var(&cfg.requestTimeout, "request_timeout", "request timeout")
	flagSet.Var(&cfg.keepAlive, "keepalive", "keepalive time")
	flagSet.BoolVarP(&cfg.Verbose, "verbose", "v", envBool(verboseKey, "true"), "verbose")
	flagSet.Uint32Var(&cfg.ringBuffer, "ring_buffer", envUint32(ringbufferKey, 512), "ring buffer size")
	flagSet.StringVar(&cfg.specificPrefix, "prefix", env(specificPrefixKey, "mesos"), "specific prefix which container name that not killed by this app")
}
