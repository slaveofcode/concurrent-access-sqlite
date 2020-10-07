package process

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/slaveofcode/concurrent-access-sqlite/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// ProcessCfg is struct for spawning new process
type ProcessCfg struct {
	PID      int
	DBPath   string
	Label    string
	DelayOps time.Duration
}

var label string

func print(text ...interface{}) {
	// log.Print("[" + label + "] ")
	// log.Print(text...)
	// log.Printf("\n")
}

func connDB(dbPath string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})

	if err != nil {
		panic("[Process One] Unable to open database connection: " + err.Error())
	}

	return db
}

func readProduct(db *gorm.DB, name string) *models.Product {
	var pdt models.Product
	db.Where("name = ?", name).First(&pdt)
	return &pdt
}

func writeProduct(db *gorm.DB, p *models.Product) {
	print("Will create new Product" + p.Name)

	err := db.Create(p).Error
	if err != nil {
		panic("Got error when creating new product" + p.Name + ": " + err.Error())
	}
}

func updateProduct(db *gorm.DB, p *models.Product) {
	err := db.Model(p).Update("name", p.Name+" - Updated").Error
	if err != nil {
		panic("Got error when updating product" + p.Name + ": " + err.Error())
	}
}

func randIsAvailable() bool {
	rand.Seed(time.Now().UnixNano())
	if rand.Intn(2) == 1 {
		return true
	} else {
		return false
	}
}

// Serve will spawn a new process
func Serve(cfg *ProcessCfg) {
	print("Process " + cfg.Label + " is running...")
	db := connDB(cfg.DBPath)

	label = cfg.Label

	for {
		time.Sleep(cfg.DelayOps * time.Millisecond)

		nameID := rand.Intn(10000-2) + 2
		price := rand.Intn(500000-200000) + 200000

		name := "KFC " + strconv.Itoa(nameID)

		pdt := readProduct(db, name)
		if pdt.Name != "" {
			updateProduct(db, pdt)
		}

		writeProduct(db, &models.Product{
			Name:        name,
			IsAvailable: randIsAvailable(),
			Price:       float64(price),
		})
	}
}
