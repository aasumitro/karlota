package config

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var configOnce sync.Once
var instance *Config

type (
	Config struct {
		sync.RWMutex
		ctx context.Context
		// SERVER CONFIGURATION
		ServerName    string `mapstructure:"SERVER_NAME"`
		ServerPort    string `mapstructure:"SERVER_PORT"`
		ServerDebug   bool   `mapstructure:"SERVER_DEBUG"`
		ServerVersion string `mapstructure:"SERVER_VERSION"`
		// DATABASE CONFIGURATION
		SQLiteDsnURL string `mapstructure:"SQLITE_DSN_URL"`
		// RABBITMQ CONFIGURATION
		RabbitMQDsnURL string `mapstructure:"RABBITMQ_DSN_URL"`
		// MAILER CONFIGURATION
		MailerHost     string `mapstructure:"MAILER_HOST"`
		MailerPort     int    `mapstructure:"MAILER_PORT"`
		MailerUsername string `mapstructure:"MAILER_USERNAME"`
		MailerPassword string `mapstructure:"MAILER_PASSWORD"`
		// SECRETS
		AuthSecretKey           string `mapstructure:"AUTH_SECRET_KEY"`
		VerifyAccountSecretKey  string `mapstructure:"VERIFY_ACCOUNT_SECRET_KEY"`
		ForgotPasswordSecretKey string `mapstructure:"FORGOT_PASSWORD_SECRET_KEY"`
		BaseURL                 string `mapstructure:"BASE_URL"`
		// App Infra
		Infra *Infra
	}

	Infra struct {
		// GormPool is a global variable that hold database and orm connection
		GormPool *gorm.DB
		SQLPool  *sql.DB
		// GinEngine is a global variable that hold gin engine
		GinEngine *gin.Engine
		// RabbitMQPublisher & RabbitMQConsumer is a global variable that hold amqp connection
		RabbitMQPublisher *amqp.Connection
		RabbitMQConsumer  *amqp.Connection
	}

	Option func(cfg *Config)
)

func Load() *Config {
	// notify that app try to load config file
	log.Println("Load configuration file . . . .")
	configOnce.Do(func() {
		// find environment file
		viper.AutomaticEnv()
		// error handling for specific case
		if err := viper.ReadInConfig(); err != nil {
			var configFileNotFoundError viper.ConfigFileNotFoundError
			if errors.As(err, &configFileNotFoundError) {
				// Config file not found; ignore error if desired
				log.Fatal(".env file not found!, please copy .env.example and paste as .env")
			}
			log.Fatalf("ENV_ERROR: %s", err.Error())
		}
		// notify that config file is ready
		log.Println("configuration file: ready")
		// extract config to struct
		if err := viper.Unmarshal(&instance); err != nil {
			log.Fatalf("ENV_ERROR: %s", err.Error())
		}
	})
	return instance
}

func (c *Config) With(ctx context.Context, options ...Option) *Config {
	c.ctx = ctx
	c.Infra = &Infra{}
	for _, opt := range options {
		opt(c)
	}
	return c
}

func GetServerInfo() string {
	return fmt.Sprintf("%s %s",
		instance.ServerName, instance.ServerVersion)
}

func GetServerAddr() string {
	return instance.ServerPort
}

func GetMailHost() string {
	return instance.MailerHost
}

func GetMailPort() int {
	return instance.MailerPort
}

func GetMailUsername() string {
	return instance.MailerUsername
}

func GetMailPassword() string {
	return instance.MailerPassword
}

func GetAuthSecret() string {
	return instance.AuthSecretKey
}

func GetForgotPwdSecret() string {
	return instance.ForgotPasswordSecretKey
}

func GetVerifyAccountSecret() string {
	return instance.VerifyAccountSecretKey
}

func GetBaseURL() string {
	return instance.BaseURL
}
