package app

import (
	"os"
	"time"

	"github.com/joho/godotenv"
	"grest.dev/grest"
)

func Config() {
	if !isConfigConfigured {
		config{}.configure()
		isConfigConfigured = true
	}
}

var isConfigConfigured = false

// default config
var (
	APP_VERSION = "23.03.161330"

	APP_ENV  = "local"
	APP_PORT = "4001"
	APP_URL  = "http://localhost:4001"

	IS_MAIN_SERVER = true // set to true to run migration, seed and task scheduling

	IS_GENERATE_OPEN_API_DOC = false

	// for testing
	ENV_FILE            = ""
	IS_USE_MOCK_SERVICE = false
	IS_USE_MOCK_DB      = false

	LOG_CONSOLE_ENABLED     = true           // print log to the terminal
	LOG_FILE_ENABLED        = true           // log to a file. the fields below can be skipped if this value is false
	LOG_FILE_USE_LOCAL_TIME = true           // if false log rotation filename will be use UTC time
	LOG_FILE_FILENAME       = "logs/api.log" //
	LOG_FILE_MAX_SIZE       = 100            // MB
	LOG_FILE_MAX_AGE        = 7              // days
	LOG_FILE_MAX_BACKUPS    = 0              // files

	JWT_KEY     = "f4cac8b77a8d4cb5881fac72388bb226"
	CRYPTO_KEY  = "wAGyTpFQX5uKV3JInABXXEdpgFkQLPTf"
	CRYPTO_SALT = "0de0cda7d2dd4937a1c4f7ddc43c580f"
	CRYPTO_INFO = "info"

	DB_DRIVER            = "sqlite"
	DB_HOST              = "127.0.0.1"
	DB_HOST_READ         = ""
	DB_PORT              = 5432
	DB_DATABASE          = "data.db"
	DB_USERNAME          = "sqlite"
	DB_PASSWORD          = "secret"
	DB_MAX_OPEN_CONNS    = 0
	DB_MAX_IDLE_CONNS    = 5
	DB_CONN_MAX_LIFETIME = time.Hour // on .env = "1h". Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
	DB_IS_DEBUG          = false

	REDIS_HOST      = "127.0.0.1"
	REDIS_PORT      = "6379"
	REDIS_CACHE_DB  = 3
	REDIS_REPORT_DB = 3
	REDIS_USERNAME  = ""
	REDIS_PASSWORD  = ""

	FS_DRIVER          = "local"
	FS_LOCAL_DIR_PATH  = "storages"
	FS_PUBLIC_DIR_PATH = "storages"
	FS_END_POINT       = "s3.amazonaws.com"
	FS_PORT            = 443
	FS_REGION          = "ap-southeast-3"
	FS_BUCKET_NAME     = "attachments"
	FS_ACCESS_KEY      = ""
	FS_SECRET_KEY      = ""

	TELEGRAM_ALERT_TOKEN   = ""
	TELEGRAM_ALERT_USER_ID = ""
)

type config struct{}

func (config) configure() {

	// set ENV_FILE with absolute path for the .env file to run test with .env
	envFile := os.Getenv("ENV_FILE")
	if envFile != "" {
		godotenv.Load(envFile)
	} else {
		godotenv.Load()
	}

	grest.LoadEnv("APP_ENV", &APP_ENV)
	grest.LoadEnv("APP_PORT", &APP_PORT)

	grest.LoadEnv("IS_MAIN_SERVER", &IS_MAIN_SERVER)

	grest.LoadEnv("IS_GENERATE_OPEN_API_DOC", &IS_GENERATE_OPEN_API_DOC)

	grest.LoadEnv("ENV_FILE", &ENV_FILE)
	grest.LoadEnv("IS_USE_MOCK_SERVICE", &IS_USE_MOCK_SERVICE)
	grest.LoadEnv("IS_USE_MOCK_DB", &IS_USE_MOCK_DB)

	grest.LoadEnv("LOG_CONSOLE_ENABLED", &LOG_CONSOLE_ENABLED)
	grest.LoadEnv("LOG_FILE_ENABLED", &LOG_FILE_ENABLED)
	grest.LoadEnv("LOG_FILE_USE_LOCAL_TIME", &LOG_FILE_USE_LOCAL_TIME)
	grest.LoadEnv("LOG_FILE_FILENAME", &LOG_FILE_FILENAME)
	grest.LoadEnv("LOG_FILE_MAX_SIZE", &LOG_FILE_MAX_SIZE)
	grest.LoadEnv("LOG_FILE_MAX_AGE", &LOG_FILE_MAX_AGE)
	grest.LoadEnv("LOG_FILE_MAX_BACKUPS", &LOG_FILE_MAX_BACKUPS)

	grest.LoadEnv("JWT_KEY", &JWT_KEY)
	grest.LoadEnv("CRYPTO_KEY", &CRYPTO_KEY)
	grest.LoadEnv("CRYPTO_SALT", &CRYPTO_SALT)
	grest.LoadEnv("CRYPTO_INFO", &CRYPTO_INFO)

	grest.LoadEnv("DB_DRIVER", &DB_DRIVER)
	grest.LoadEnv("DB_HOST", &DB_HOST)
	grest.LoadEnv("DB_HOST_READ", &DB_HOST_READ)
	grest.LoadEnv("DB_PORT", &DB_PORT)
	grest.LoadEnv("DB_DATABASE", &DB_DATABASE)
	grest.LoadEnv("DB_USERNAME", &DB_USERNAME)
	grest.LoadEnv("DB_PASSWORD", &DB_PASSWORD)
	grest.LoadEnv("DB_MAX_OPEN_CONNS", &DB_MAX_OPEN_CONNS)
	grest.LoadEnv("DB_MAX_IDLE_CONNS", &DB_MAX_IDLE_CONNS)
	grest.LoadEnv("DB_CONN_MAX_LIFETIME", &DB_CONN_MAX_LIFETIME)
	grest.LoadEnv("DB_IS_DEBUG", &DB_IS_DEBUG)

	grest.LoadEnv("REDIS_HOST", &REDIS_HOST)
	grest.LoadEnv("REDIS_PORT", &REDIS_PORT)
	grest.LoadEnv("REDIS_CACHE_DB", &REDIS_CACHE_DB)
	grest.LoadEnv("REDIS_REPORT_DB", &REDIS_REPORT_DB)
	grest.LoadEnv("REDIS_USERNAME", &REDIS_USERNAME)
	grest.LoadEnv("REDIS_PASSWORD", &REDIS_PASSWORD)

	grest.LoadEnv("FS_END_POINT", &FS_END_POINT)
	grest.LoadEnv("FS_PORT", &FS_PORT)
	grest.LoadEnv("FS_REGION", &FS_REGION)
	grest.LoadEnv("FS_BUCKET_NAME", &FS_BUCKET_NAME)
	grest.LoadEnv("FS_ACCESS_KEY", &FS_ACCESS_KEY)
	grest.LoadEnv("FS_SECRET_KEY", &FS_SECRET_KEY)

	grest.LoadEnv("TELEGRAM_ALERT_TOKEN", &TELEGRAM_ALERT_TOKEN)
	grest.LoadEnv("TELEGRAM_ALERT_USER_ID", &TELEGRAM_ALERT_USER_ID)
}
