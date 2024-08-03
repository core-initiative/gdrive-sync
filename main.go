package main

import (
    "flag"
    "fmt"
    "log"
    "github.com/robfig/cron/v3"
)

func main() {
    // Define command-line flags
    configPath := flag.String("config", "config.yaml", "Path to the configuration file")
    immediate := flag.Bool("immediate", false, "Run the backup immediately without waiting for the schedule")
    flag.Parse()

    config, err := LoadConfig(*configPath)
    if err != nil {
        log.Fatalf("Error loading configuration: %v", err)
    }

    if *immediate {
        fmt.Printf("Running backup immediately for folder: %s\n", config.Folder)
        SyncFolder(config.Folder, config.TargetFolderID)
        return
    }

    c := cron.New()
    _, err = c.AddFunc(config.Schedule, func() {
        fmt.Printf("Syncing folder: %s\n", config.Folder)
        SyncFolder(config.Folder, config.TargetFolderID)
    })

    if err != nil {
        log.Fatalf("Error scheduling sync: %v", err)
    }

    c.Start()
    defer c.Stop()

    // Keep the application running
    select {}
}