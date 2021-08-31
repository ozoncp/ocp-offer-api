package config

import (
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

// Build information -ldflags .
var (
	version    string = "dev"
	commitHash string = "-"
)

var cfg *config

var (
	Project  *project
	GRPC     *gRPC
	Gateway  *gateway
	Database *database
	Metrics  *metrics
	Kafka    *kafka
	Status   *status
)

// config - microservice config.
type config struct {
	Project  project  `yaml:"project"`
	GRPC     gRPC     `yaml:"grpc"`
	Gateway  gateway  `yaml:"gateway"`
	Metrics  metrics  `yaml:"metrics"`
	Database database `yaml:"database"`
	Kafka    kafka    `yaml:"kafka"`
	Status   status   `yaml:"status"`
}

// gRPC config.
type project struct {
	Debug       bool   `yaml:"debug"`
	Name        string `yaml:"name"`
	Environment string `yaml:"environment"`
	Version     string
	CommitHash  string
}

// gRPC config.
type gRPC struct {
	Port              int    `yaml:"port" env:"GRPC_PORT"`
	MaxConnectionIdle int64  `yaml:"maxConnectionIdle" env:"GRPC_MAX_CONN_IDLE"`
	Timeout           int64  `yaml:"timeout" env:"GRPC_TIMEOUT"`
	MaxConnectionAge  int64  `yaml:"maxConnectionAge" env:"GRPC_CONN_AGE"`
	Host              string `yaml:"host" env:"GRPC_HOST"`
}

// gateway config.
type gateway struct {
	Port int    `yaml:"port" env:"GATEWAY_PORT"`
	Host string `yaml:"host" env:"GATEWAY_HOST"`
}

type metrics struct {
	Port int    `yaml:"port" env:"METRICS_PORT"`
	Host string `yaml:"host" env:"METRICS_HOST"`
	Path string `yaml:"path" env:"METRICS_PATH"`
}

// Postgres config.
type database struct {
	Port     int    `yaml:"port" env:"DATABASE_PORT"`
	Host     string `yaml:"host" env:"DATABASE_HOST"`
	User     string `yaml:"user" env:"DATABASE_USER"`
	Password string `yaml:"password" env:"DATABASE_PASSWORD"`
	Name     string `yaml:"name" env:"DATABASE_NAME"`
	SSLMode  string `yaml:"sslMode" env:"DATABASE_SSL_MODE"`
	Driver   string `yaml:"driver" env:"DATABASE_DRIVER"`
}

// Kafka config.
type kafka struct {
	Capacity uint64   `yaml:"capacity"`
	Topic    string   `yaml:"topic"`
	GroupID  string   `yaml:"groupId"`
	Brokers  []string `yaml:"brokers"`
}

// Service status config.
type status struct {
	Port          int    `yaml:"port" env:"STATUS_PORT"`
	Host          string `yaml:"host" env:"STATUS_HOST"`
	VersionPath   string `yaml:"versionPath" env:"STATUS_VERSION_PATH"`
	LivenessPath  string `yaml:"livenessPath" env:"STATUS_LIVENESS_PATH"`
	ReadinessPath string `yaml:"readinessPath" env:"STATUS_READINESS_PATH"`
}

var fileConfig = "config.yml"
var doOnce sync.Once

func init() {
	// Once initialized
	doOnce.Do(func() {
		if err := UpdateConfig(); err != nil {
			log.Fatal().Err(err).Msg("Configuration initialization failed")

			return
		}

		hotReload := make(chan os.Signal, 1)
		signal.Notify(hotReload, syscall.SIGHUP)

		go func() {
			for {
				<-hotReload
				log.Info().Msgf("Hot reload configurations ...")
				if err := UpdateConfig(); err != nil {
					log.Error().Err(err).Msgf("Error on reloading config")
					os.Exit(1)
				}
				log.Info().Msg("Config hot reloading was successful")
			}
		}()

		log.Info().Msg("Config initialization was successful")
	})
}

// UpdateConfig - Updates the config by rereading.
func UpdateConfig() error {
	file, err := os.Open(fileConfig)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return err
	}

	// read environment and replace
	readEnvAndSet(reflect.ValueOf(cfg))

	// Set global value
	Project = &cfg.Project
	Project.Version = version
	Project.CommitHash = commitHash

	GRPC = &cfg.GRPC
	Gateway = &cfg.Gateway
	Metrics = &cfg.Metrics
	Database = &cfg.Database
	Kafka = &cfg.Kafka
	Status = &cfg.Status

	return nil
}

// readEnvAndSet - Sets config from environment values.
func readEnvAndSet(v reflect.Value) {
	if v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}

	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Type.Kind() == reflect.Struct {
			readEnvAndSet(v.Field(i))
		} else if tag := field.Tag.Get("env"); tag != "" {
			if value := os.Getenv(tag); value != "" {
				if err := setValue(v.Field(i), value); err != nil {
					log.Error().Err(err).Msgf("Failed to set environment value for \"%s\"", field.Name)
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

	// unsupported types
	case reflect.Map, reflect.Ptr,
		reflect.Complex64, reflect.Interface,
		reflect.Invalid, reflect.Slice, reflect.Func,
		reflect.Array, reflect.Chan, reflect.Complex128,
		reflect.Struct, reflect.Uintptr, reflect.UnsafePointer:
	default:
		return fmt.Errorf("unsupported type: %v", valueType.Kind())
	}

	return nil
}
