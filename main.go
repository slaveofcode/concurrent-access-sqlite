package main

import (
	"log"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"syscall"

	"github.com/slaveofcode/concurrent-access-sqlite/models"
	"github.com/slaveofcode/concurrent-access-sqlite/process"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func getCurrentWorkingPath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

func getFullPathDB() string {
	return path.Join(getCurrentWorkingPath(), "app.db")
}

func createFreshDB(refreshDB bool) *gorm.DB {
	dbPath := getFullPathDB()

	if refreshDB {
		_, err := os.Stat(dbPath)
		if !os.IsNotExist(err) {
			log.Println("Removing old DB...")
			if err = os.Remove(dbPath); err != nil {
				panic("Unable to remove old DB")
			}
		}

		log.Println("Create new DB at:", getFullPathDB())
	}

	// busy_timeout 5 secs
	db, err := gorm.Open(sqlite.Open(dbPath+"?busy_timeout=5000"), &gorm.Config{})

	db.Exec("PRAGMA journal_mode=WAL")

	if err != nil {
		panic("Unable to open database connection: " + err.Error())
	}

	db.AutoMigrate(&models.Product{})

	return db
}

func main() {
	db := createFreshDB(false)
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	dbPath := getFullPathDB()
	processes := []process.ProcessCfg{
		{
			Label:    "Process One",
			PID:      1,
			DBPath:   dbPath,
			DelayOps: 1,
		},
		{
			Label:    "Process Two",
			PID:      2,
			DBPath:   dbPath,
			DelayOps: 1,
		},
		{
			Label:    "Process Three",
			PID:      3,
			DBPath:   dbPath,
			DelayOps: 0,
		},
		{
			Label:    "Process Four",
			PID:      4,
			DBPath:   dbPath,
			DelayOps: 0,
		},
		{
			Label:    "Process Five",
			PID:      5,
			DBPath:   dbPath,
			DelayOps: 0,
		},
		{
			Label:    "Process Six",
			PID:      6,
			DBPath:   dbPath,
			DelayOps: 0,
		},
		{
			Label:    "Process Seven",
			PID:      7,
			DBPath:   dbPath,
			DelayOps: 0,
		},
		{
			Label:    "Process Eight",
			PID:      8,
			DBPath:   dbPath,
			DelayOps: 0,
		},
	}

	for _, proc := range processes {
		log.Println("Spawning " + proc.Label)
		go process.Serve(&proc)
	}

	exitChan := make(chan os.Signal, 2)
	signal.Notify(exitChan, syscall.SIGINT)

	for {
		select {
		case <-exitChan:
			log.Println("Shutting down...")
			os.Exit(0)
		}
	}
}
