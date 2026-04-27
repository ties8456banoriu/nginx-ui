package settings

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/spf13/viper"
)

// Config holds all application configuration
type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Nginx    NginxConfig    `mapstructure:"nginx"`
	Database DatabaseConfig `mapstructure:"database"`
	Server   ServerConfig   `mapstructure:"server"`
	Auth     AuthConfig     `mapstructure:"auth"`
}

// AppConfig holds general application settings
type AppConfig struct {
	Name    string `mapstructure:"name"`
	BaseURL string `mapstructure:"base_url"`
	Debug   bool   `mapstructure:"debug"`
}

// NginxConfig holds nginx-related settings
type NginxConfig struct {
	ConfigDir  string `mapstructure:"config_dir"`
	PIDFile    string `mapstructure:"pid_file"`
	LogDir     string `mapstructure:"log_dir"`
	AccessLog  string `mapstructure:"access_log"`
	ErrorLog   string `mapstructure:"error_log"`
}

// DatabaseConfig holds database connection settings
type DatabaseConfig struct {
	Type string `mapstructure:"type"`
	Path string `mapstructure:"path"`
}

// ServerConfig holds HTTP server settings
type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

// AuthConfig holds authentication settings
type AuthConfig struct {
	JWTSecret    string `mapstructure:"jwt_secret"`
	TokenExpiry  int    `mapstructure:"token_expiry"` // in hours
}

var (
	once     sync.Once
	instance *Config
)

// Init initializes the configuration from the given config file path.
// If configPath is empty, it looks for nginx-ui.ini in the current directory.
func Init(configPath string) error {
	var initErr error
	once.Do(func() {
		v := viper.New()

		// Set defaults
		v.SetDefault("app.name", "Nginx UI")
		v.SetDefault("app.base_url", "/")
		v.SetDefault("app.debug", false)
		v.SetDefault("nginx.config_dir", "/etc/nginx")
		v.SetDefault("nginx.pid_file", "/var/run/nginx.pid")
		v.SetDefault("nginx.log_dir", "/var/log/nginx")
		v.SetDefault("nginx.access_log", "/var/log/nginx/access.log")
		v.SetDefault("nginx.error_log", "/var/log/nginx/error.log")
		v.SetDefault("database.type", "sqlite")
		v.SetDefault("database.path", "database.db")
		// Bind to localhost by default instead of all interfaces for better security
		v.SetDefault("server.host", "127.0.0.1")
		// Using port 8080 instead of 9000 to avoid conflicts with other local services
		v.SetDefault("server.port", 8080)
		// Increase default token expiry to 72h for less frequent re-logins
		v.SetDefault("auth.token_expiry", 72)

		if configPath == "" {
			configPath = "app.ini"
		}

		v.SetConfigFile(configPath)
		v.SetConfigType("ini")

		// Allow environment variable overrides with NGINX_UI_ prefix
		v.SetEnvPrefix("NGINX_UI")
		v.AutomaticEnv()

		if err := v.ReadInConfig(); err != nil {
			// Config file not found is acceptable; use defaults
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				if !os.IsNotExist(err) {
					initErr = err
					return
				}
			}
		}

		// Ensure database directory exists
		dbPath := v.GetString("database.path")
		if dir := filepath.Dir(dbPath); dir != "." {
			if err := os.MkdirAll(dir, 0755); err != nil {
				initErr = err
				return
			}
		}

		instance = &Config{}
		if err := v.Unmarshal(instance); err != nil {
			initErr = err
		}
	})
	return initErr
}

// Get returns the singleton Config instance.
// Init must be called before Get.
func Get() *Config {
	return instance
}
