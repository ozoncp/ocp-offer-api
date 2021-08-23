package config

import (
	"os"
	"reflect"
	"strconv"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

var cfg *config

var (
	Project  *project
	GRPC     *gRPC
	Gateway  *gateway
	Database *database
	Metrics  *metrics
	Kafka    *kafka
)

// config - microservice config
type config struct {
	Project  project  `yaml:"Project"`
	GRPC     gRPC     `yaml:"GRPC"`
	Gateway  gateway  `yaml:"Gateway"`
	Metrics  metrics  `yaml:"Metrics"`
	Database database `yaml:"Database"`
	Kafka    kafka    `yaml:"Kafka"`
}

// gRPC config
type project struct {
	Name        string `yaml:"Name"`
	Version     string `yaml:"Version"`
	Environment string `yaml:"Environment"`
	Debug       bool   `yaml:"Debug"`
}

// gRPC config
type gRPC struct {
	Host              string `yaml:"Host" env:"GRPC_HOST"`
	Port              int    `yaml:"Port" env:"GRPC_PORT"`
	MaxConnectionIdle int64  `yaml:"MaxConnectionIdle" env:"GRPC_MAX_CONN_IDLE"`
	Timeout           int64  `yaml:"Timeout" env:"GRPC_TIMEOUT"`
	MaxConnectionAge  int64  `yaml:"MaxConnectionAge" env:"GRPC_CONN_AGE"`
}

// gateway config
type gateway struct {
	Host string `yaml:"Host" env:"GATEWAY_HOST"`
	Port int    `yaml:"Port" env:"GATEWAY_PORT"`
}

type metrics struct {
	Host string `yaml:"Host" env:"METRICS_HOST"`
	Port int    `yaml:"Port" env:"METRICS_PORT"`
	Path string `yaml:"Path" env:"METRICS_PATH"`
}

// Postgres config
type database struct {
	Host     string `yaml:"Host" env:"DATABASE_HOST"`
	Port     int    `yaml:"Port" env:"DATABASE_PORT"`
	User     string `yaml:"User" env:"DATABASE_USER"`
	Password string `yaml:"Password" env:"DATABASE_PASSWORD"`
	Name     string `yaml:"Name" env:"DATABASE_NAME"`
	SSLMode  string `yaml:"SslMode" env:"DATABASE_SSL_MODE"`
	Driver   string `yaml:"Driver" env:"DATABASE_DRIVER"`
}

// Postgres config
type kafka struct {
	Brokers  []string `yaml:"Brokers"`
	Topic    string   `yaml:"Topic"`
	Capacity uint64   `yaml:"Capacity"`
}

var fileConfig = "config.yml"
var doOnce sync.Once

func init() {
	// Once initialized
	doOnce.Do(func() {
		if err := UpdateConfig(); err != nil {
			log.Fatal().Err(err).Msg("Configuration initialization failed")
		} else {
			log.Info().Msg("Config initialization was successful")
			return
		}
	})
}

// UpdateConfig - Updates the config by rereading
func UpdateConfig() error {
	file, err := os.Open(fileConfig)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err = decoder.Decode(&cfg); err != nil {
		return err
	}

	// read environment and replace
	readEnvAndSet(reflect.ValueOf(cfg))

	// Set global value
	Project = &cfg.Project
	GRPC = &cfg.GRPC
	Gateway = &cfg.Gateway
	Metrics = &cfg.Metrics
	Database = &cfg.Database
	Kafka = &cfg.Kafka

	return nil
}

// readEnvAndSet - Sets config from environment values
func readEnvAndSet(v reflect.Value) {
	if v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}

	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Type.Kind() == reflect.Struct {
			readEnvAndSet(v.Field(i))
		} else {
			if tag := field.Tag.Get("env"); tag != "" {
				if value := os.Getenv(tag); value != "" {
					if err := setValue(v.Field(i), value); err != nil {
						log.Error().Err(err).Msgf("Failed to set environment value for \"%s\"", field.Name)
					}
				}
			}
		}
	}
}

func setValue(field reflect.Value, value string) error {
	valueType := field.Type()

	switch valueType.Kind() {
	// set string value
	case reflect.String:
		field.SetString(value)

	// set boolean value
	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		field.SetBool(b)

	// set integer (or time) value
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if field.Kind() == reflect.Int64 && valueType.PkgPath() == "time" && valueType.Name() == "Duration" {
			// try to parse time
			d, err := time.ParseDuration(value)
			if err != nil {
				return err
			}
			field.SetInt(int64(d))

		} else {
			// parse regular integer
			number, err := strconv.ParseInt(value, 0, valueType.Bits())
			if err != nil {
				return err
			}
			field.SetInt(number)
		}

	// set unsigned integer value
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		number, err := strconv.ParseUint(value, 0, valueType.Bits())
		if err != nil {
			return err
		}
		field.SetUint(number)

	// set floating point value
	case reflect.Float32, reflect.Float64:
		number, err := strconv.ParseFloat(value, valueType.Bits())
		if err != nil {
			return err
		}
		field.SetFloat(number)
	}

	return nil
}
