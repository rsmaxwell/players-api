package config

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/rsmaxwell/players-api/internal/debug"
)

// Database type
type Database struct {
	DriverName   string `json:"driverName"`
	UserName     string `json:"userName"`
	Password     string `json:"password"`
	Scheme       string `json:"scheme"`
	Host         string `json:"host"`
	Path         string `json:"path"`
	DatabaseName string `json:"databaseName"`
}

// Server type
type Server struct {
	Port int `json:"port"`
}

// Config type
type ConfigFile struct {
	Database           Database `json:"database"`
	Server             Server   `json:"server"`
	AccessTokenExpiry  string   `json:"accessToken_expiry"`
	RefreshTokenExpiry string   `json:"refreshToken_expiry"`
	ClientRefreshDelta string   `json:"clientRefreshDelta"`
}

// Config type
type Config struct {
	Database           Database
	Server             Server
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
	ClientRefreshDelta time.Duration
}

var (
	pkg           = debug.NewPackage("config")
	functionOpen  = debug.NewFunction(pkg, "Open")
	functionSetup = debug.NewFunction(pkg, "Setup")
)

// Setup function
func Setup() (*sql.DB, *Config, error) {
	f := functionSetup

	// Read configuration
	rootDir := debug.RootDir()
	c, err := Open(rootDir)
	if err != nil {
		message := "Could not open configuration"
		f.Errorf(message)
		f.DumpError(err, message)
		return nil, nil, err
	}

	// Connect to the database
	db, err := sql.Open(c.DriverName(), c.ConnectionString())
	if err != nil {
		message := "Could not connect to the database"
		f.Errorf(message)
		f.DumpError(err, message)
		return nil, nil, err
	}

	return db, c, nil
}

// Open returns the configuration
func Open(rootDir string) (*Config, error) {
	f := functionOpen

	configFileName := filepath.Join(rootDir, "config", "config.json")

	bytearray, err := ioutil.ReadFile(configFileName)
	if err != nil {
		f.Dump("could not read config file")
		return nil, err
	}

	var configFile ConfigFile
	err = json.Unmarshal(bytearray, &configFile)
	if err != nil {
		f.Dump("could not Unmarshal configuration")
		return nil, err
	}

	return configFile.toConfig()
}
