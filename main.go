package main

import (
	"context"
	"flag"

    "drive-sync/log"
	"github.com/robfig/cron/v3"
)
var Version = "0.1"
func main() {

    // Define command-line flags
    configPath := flag.String("config", "./conf/config.yaml", "Path to the configuration file")
    immediate := flag.Bool("immediate", false, "Run the backup immediately without waiting for the schedule")
    flag.Parse()

    // create root logger tagged with server version
	logger := log.New().With(context.Background(), "version", Version)

    config, err := LoadConfig(*configPath, logger)
    if err != nil {
        logger.Errorf("Error loading configuration: %v", err)
    }

    if *immediate {
        logger.Infof("Running backup immediately for folder: %s\n", config.Folder)
        SyncFolder(config.Folder, config.TargetFolderID)
        return
    }

    c := cron.New()
    _, err = c.AddFunc(config.Schedule, func() {
        logger.Infof("Syncing folder: %s\n", config.Folder)
        SyncFolder(config.Folder, config.TargetFolderID)
    })

    if err != nil {
        logger.Errorf("Error scheduling sync: %v", err)
    }

    c.Start()
    defer c.Stop()

    // Keep the application running
    select {}
}