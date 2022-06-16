package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kramllih/filterService/api"
	"github.com/kramllih/filterService/config"
	"github.com/spf13/viper"

	"github.com/kramllih/filterService/internal/database"
	_ "github.com/kramllih/filterService/internal/database/bbolt"
	_ "github.com/kramllih/filterService/internal/logger"
)

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
}

func main() {

	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

}

func run() error {

	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("unable to read config: %w", err)
	}

	var (
		c      *config.RawConfig
		apicfg *config.RawConfig
		ls     string
	)

	err = viper.Unmarshal(&c)
	if err != nil {
		return fmt.Errorf("config unmarshal error: %w", err)
	}

	databaseCfg, err := config.UnpackNamespace("database", c)
	if err != nil {
		return err
	}

	db, err := database.Load(&databaseCfg)
	if err != nil {
		return fmt.Errorf("error loading database: %w", err)
	}

	path, err := getWorkingPath()
	if err != nil {
		return err
	}

	err = c.UnpackAttribute("api", &apicfg)
	if err != nil {
		return err
	}

	err = c.UnpackAttribute("languageservice", &ls)
	if err != nil {
		return err
	}

	router, err := api.SetupRouter(path, db, apicfg, ls)
	if err != nil {
		return err
	}

	return router.Start()

}

// getWorkingPath gets the working path of the application
func getWorkingPath() (string, error) {

	fullexecpath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %w", err)
	}

	dir, _ := filepath.Split(fullexecpath)

	return dir, nil

}
