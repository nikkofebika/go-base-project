package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	AppName    string
	AppPort    int
	AppEnv     string
	AppPrefork bool
	WebUrl     string

	AwsAccessKeyID     string
	AwsSecretAccessKey string
	AwsDefaultRegion   string
	AwsBucket          string

	CorsOrigins []string

	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string

	FilesystemDisk string

	JWTSecret           string
	JWTRefreshSecret    string
	JWTExpiresIn        time.Duration
	JWTRefreshExpiresIn time.Duration

	MailFrom string
	MailHost string
	MailPort int
	MailUser string
	MailPass string

	LogFormat string
	LogFile   string

	MediaDisk          string
	MediaPublicBaseUrl string

	PasswordBypass string
}

func NewAppConfig() *AppConfig {
	if os.Getenv("APP_ENV") != "production" {
		if err := godotenv.Load(".env"); err != nil {
			panic(err)
		}
	}

	return &AppConfig{
		AppName:    getString("APP_NAME", "Go Fiber MYSQL"),
		AppPort:    getInt("APP_PORT", 8080),
		AppEnv:     getString("APP_ENV", "development"),
		WebUrl:     getString("WEB_URL", "http://localhost:8000"),
		AppPrefork: getBool("APP_PREFORK", true),

		AwsAccessKeyID:     getString("AWS_ACCESS_KEY_ID", ""),
		AwsSecretAccessKey: getString("AWS_SECRET_ACCESS_KEY", ""),
		AwsDefaultRegion:   getString("AWS_DEFAULT_REGION", ""),
		AwsBucket:          getString("AWS_BUCKET", ""),

		CorsOrigins: func() []string {
			v := os.Getenv("CORS_ORIGINS")
			if v == "" {
				return []string{}
			}
			return strings.Split(v, ",")
		}(),

		DBHost:     getString("DB_HOST", "localhost"),
		DBPort:     getString("DB_PORT", "3306"),
		DBName:     getString("DB_NAME", "ticket_support"),
		DBUser:     getString("DB_USER", "root"),
		DBPassword: getString("DB_PASSWORD", ""),

		FilesystemDisk: getString("FILESYSTEM_DISK", "local"),

		JWTSecret:           mustString("JWT_SECRET"),
		JWTRefreshSecret:    mustString("JWT_REFRESH_SECRET"),
		JWTExpiresIn:        getDuration("JWT_EXPIRES_IN", time.Hour*24),
		JWTRefreshExpiresIn: getDuration("JWT_REFRESH_EXPIRES_IN", time.Hour*24*30),

		MailFrom: getString("MAIL_FROM", "info@app.com"),
		MailHost: getString("MAIL_HOST", "sandbox.smtp.mailtrap.io"),
		MailPort: getInt("MAIL_PORT", 2525),
		MailUser: getString("MAIL_USERNAME", ""),
		MailPass: getString("MAIL_PASSWORD", ""),

		LogFormat: getString("LOG_FORMAT", "json"),
		LogFile:   getString("LOG_FILE", "./storage/logs/app.log"),

		MediaDisk:          getString("MEDIA_DISK", "local"),
		MediaPublicBaseUrl: getString("MEDIA_PUBLIC_BASE_URL", "http://localhost:8080/assets/"),

		PasswordBypass: getString("PASSWORD_BYPASS", "passwordbypass"),
	}
}

func getString(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

func mustString(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("environment variable not set: %s", key)
	}
	return v
}

func getInt(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}

	i, err := strconv.Atoi(v)
	if err != nil {
		log.Fatalf("Invalid int for %s", key)
	}
	return i
}

func getDuration(key string, def time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return def
	}

	d, err := time.ParseDuration(v)
	if err != nil {
		log.Fatalf("Invalid duration for %s", key)
	}
	return d
}

func getBool(key string, def bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return def
	}

	b, err := strconv.ParseBool(v)
	if err != nil {
		log.Fatalf("Invalid bool for %s", key)
	}
	return b
}
