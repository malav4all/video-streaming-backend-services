package config

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all application configuration values, resolved for the
// currently active environment (APP_ENV).
type Config struct {
	AppEnv   string
	LogLevel string
	AppPort  string

	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	DBSSLMode      string
	DBMaxOpenConns int
	DBMaxIdleConns int

	JWTSecret     string
	JWTExpiration time.Duration

	RateLimitPerMinute int
	MaxRequestSize     int64
	RequestTimeout     time.Duration

	AllowedOrigins       []string
	CORSAllowCredentials bool

	KeycloakBaseURL      string
	KeycloakRealm        string
	KeycloakClientID     string
	KeycloakClientSecret string
	KeycloakRedirectURL  string

	RedisHost     string
	RedisPort     string
	RedisUsername string
	RedisPassword string
	RedisDB       int
	SessionTTL    time.Duration
}

// envSuffix maps APP_ENV to the suffix used on per-environment variables,
// e.g. APP_ENV=local -> reads DB_HOST_LOCAL, JWT_SECRET_LOCAL, etc.
func envSuffix(appEnv string) string {
	switch strings.ToLower(appEnv) {
	case "local":
		return "LOCAL"
	case "development", "dev":
		return "DEV"
	case "uat":
		return "UAT"
	case "production", "prod":
		return "PROD"
	default:
		return "LOCAL"
	}
}

// Load reads configuration from .env (if present) and environment variables.
// It first determines APP_ENV, then resolves every other setting using the
// matching per-environment suffix (_LOCAL / _DEV / _UAT / _PROD).
func Load() *Config {
	_ = godotenv.Load()

	appEnv := getEnv("APP_ENV", "local")
	suf := envSuffix(appEnv)

	return &Config{
		AppEnv:   appEnv,
		LogLevel: getEnv("LOG_LEVEL", "info"),
		AppPort:  getEnv("PORT", "8080"),

		DBHost:         getEnv("DB_HOST_"+suf, "localhost"),
		DBPort:         getEnv("DB_PORT_"+suf, "5432"),
		DBUser:         getEnv("DB_USER_"+suf, "postgres"),
		DBPassword:     getEnv("DB_PASSWORD_"+suf, ""),
		DBName:         getEnv("DB_NAME_"+suf, "userdb"),
		DBSSLMode:      getEnv("DB_SSLMODE_"+suf, "disable"),
		DBMaxOpenConns: getEnvAsInt("DB_MAX_OPEN_CONNS_"+suf, 25),
		DBMaxIdleConns: getEnvAsInt("DB_MAX_IDLE_CONNS_"+suf, 10),

		JWTSecret:     getEnv("JWT_SECRET_"+suf, ""),
		JWTExpiration: getEnvAsDuration("JWT_EXPIRATION_"+suf, 24*time.Hour),

		RateLimitPerMinute: getEnvAsInt("RATE_LIMIT_PER_MINUTE_"+suf, 1000),
		MaxRequestSize:     getEnvAsInt64("MAX_REQUEST_SIZE_"+suf, 104857600),
		RequestTimeout:     getEnvAsDuration("REQUEST_TIMEOUT_"+suf, 60*time.Second),

		AllowedOrigins:       getEnvAsSlice("ALLOWED_ORIGINS_"+suf, []string{"*"}),
		CORSAllowCredentials: getEnvAsBool("CORS_ALLOW_CREDENTIALS_"+suf, false),

		KeycloakBaseURL:      getEnv("KEYCLOAK_URL_"+suf, "http://localhost:8080"),
		KeycloakRealm:        getEnv("KEYCLOAK_REALM_"+suf, "master"),
		KeycloakClientID:     getEnv("KEYCLOAK_CLIENT_ID_"+suf, ""),
		KeycloakClientSecret: getEnv("KEYCLOAK_CLIENT_SECRET_"+suf, ""),
		KeycloakRedirectURL:  getEnv("KEYCLOAK_REDIRECT_URL_"+suf, ""),

		RedisHost:     getEnv("REDIS_HOST_"+suf, "localhost"),
		RedisPort:     getEnv("REDIS_PORT_"+suf, "6379"),
		RedisUsername: getEnv("REDIS_USERNAME_"+suf, ""),
		RedisPassword: getEnv("REDIS_PASSWORD_"+suf, ""),
		RedisDB:       getEnvAsInt("REDIS_DB_"+suf, 0),
		SessionTTL:    getEnvAsDuration("SESSION_TTL_"+suf, 24*time.Hour),
	}
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	if v, ok := os.LookupEnv(key); ok {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}

func getEnvAsInt64(key string, fallback int64) int64 {
	if v, ok := os.LookupEnv(key); ok {
		if i, err := strconv.ParseInt(v, 10, 64); err == nil {
			return i
		}
	}
	return fallback
}

func getEnvAsBool(key string, fallback bool) bool {
	if v, ok := os.LookupEnv(key); ok {
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	}
	return fallback
}

func getEnvAsDuration(key string, fallback time.Duration) time.Duration {
	if v, ok := os.LookupEnv(key); ok {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return fallback
}

func getEnvAsSlice(key string, fallback []string) []string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		parts := strings.Split(v, ",")
		for i := range parts {
			parts[i] = strings.TrimSpace(parts[i])
		}
		return parts
	}
	return fallback
}
