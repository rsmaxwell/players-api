package config

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

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
type Config struct {
	Database Database `json:"database"`
	Server   Server   `json:"server"`
}

var (
	pkg           = debug.NewPackage("config")
	functionOpen  = debug.NewFunction(pkg, "Open")
	functionSetup = debug.NewFunction(pkg, "Setup")
)

// Open returns the configuration
func Open(rootDir string) (*Config, error) {
	f := functionOpen

	configFile := filepath.Join(rootDir, "config", "config.json")

	bytearray, err := ioutil.ReadFile(configFile)
	if err != nil {
		f.Dump("could not read config file")
		return nil, err
	}

	var config Config
	err = json.Unmarshal(bytearray, &config)
	if err != nil {
		f.Dump("could not Unmarshal configuration")
		return nil, err
	}

	return &config, nil
}

// DriverName returns the driver name for the configured database
func (c *Config) DriverName() string {
	return c.Database.DriverName
}

// ConnectionString returns the string used to connect to the database
func (c *Config) ConnectionString() string {
	return fmt.Sprintf("%s://%s:%s@%s/%s", c.Database.Scheme, c.Database.UserName, c.Database.Password, c.Database.Host, c.Database.DatabaseName)
}

// ConnectionStringBasic returns the string used to connect to the database
func (c *Config) ConnectionStringBasic() string {
	return fmt.Sprintf("%s://%s:%s@%s", c.Database.Scheme, c.Database.UserName, c.Database.Password, c.Database.Host)
}

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

// Setup function
func SetupBasic() (*sql.DB, *Config, error) {
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
	db, err := sql.Open(c.DriverName(), c.ConnectionStringBasic())
	if err != nil {
		message := "Could not connect to the database"
		f.Errorf(message)
		f.DumpError(err, message)
		return nil, nil, err
	}

	return db, c, nil
}
