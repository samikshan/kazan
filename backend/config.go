package backend

import (
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/samikshan/kazan/backend/logging"
)

const (
	defaultloglevel = "debug"
	defaultHTTPPort = 1323
)

type Constants struct {
	HTTPPort    uint64
	LogLevel    string
	PowHost     string
	PowGrpcHost string
	Postgres    struct {
		Type     string
		Host     string
		Port     uint16
		User     string
		Password string
		Database string
	}
	JWT struct {
		SigningKey string
	}
}

type Config struct {
	Constants
}

var (
	Cfg *Config
)

// initialise global configuration variable
func init() {
	log.Info("Initialising configuration...")
	if Cfg != nil {
		log.Info(Cfg)
		return
	}

	constants, err := initViper()
	if err != nil {
		log.WithError(err).Fatalf("failed to initialise configuration")
	}

	Cfg = &Config{
		Constants: constants,
	}
}

func initFlags() {
	flag.String("config", "", "Configuration file for LetsJam")
	flag.String(logging.LevelFlag, defaultloglevel, logging.LevelHelp)
}

func initViper() (Constants, error) {
	initFlags()
	if err := viper.BindPFlags(flag.CommandLine); err != nil {
		return Constants{}, err
	}

	viper.SetConfigName("kazan-backend.config") // Configuration fileName without the .TOML or .YAML extension
	viper.AddConfigPath(".")                    // Search the root directory for the configuration file

	if confFile := viper.GetString("config"); confFile != "" {
		viper.SetConfigFile(confFile)
	}

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	var constants Constants
	err := viper.Unmarshal(&constants)
	return constants, err
}
