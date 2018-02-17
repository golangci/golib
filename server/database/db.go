package database

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/golangci/golib/server/cfg"
	"github.com/jinzhu/gorm"
	yaml "gopkg.in/yaml.v2"
)

type DBConfig struct {
	Adapter  string
	Database string
	Host     string
	Port     int
	User     string
	Enabled  bool

	ConnString string
}

func readYamlDBConfig(filePath string) (*DBConfig, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var r DBConfig
	err = yaml.Unmarshal(content, &r)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func connectToDB(cfg DBConfig) (*gorm.DB, error) {
	log.Printf("connecting to %s", cfg.ConnString)
	db, err := gorm.Open(cfg.Adapter, cfg.ConnString)
	if err != nil {
		return nil, err
	}

	log.Printf("connected to DB")
	return db, nil
}

func connectToDBP(cfg DBConfig) *gorm.DB {
	db, err := connectToDB(cfg)
	if err != nil {
		panic(err)
	}

	return db
}

func GetDBConfig() (*DBConfig, error) {
	if dbEnvURL := os.Getenv("DATABASE_URL"); dbEnvURL != "" {
		dbEnvURL = strings.Replace(dbEnvURL, "postgresql", "postgres", 1)
		adapter := strings.Split(dbEnvURL, "://")[0]
		return &DBConfig{
			Adapter:    adapter,
			ConnString: dbEnvURL,
		}, nil
	}

	DBConfigPath := cfg.GetRoot() + "db.yaml"

	cfg, err := readYamlDBConfig(DBConfigPath)
	if err != nil {
		return nil, err
	}

	if !cfg.Enabled {
		return nil, nil
	}

	cfg.ConnString = fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Database)
	return cfg, nil
}

var db *gorm.DB
var dbOnce sync.Once

func initDB() {
	cfg, err := GetDBConfig()
	if err != nil {
		log.Fatalf("Can't get db config: %s", err)
	}

	db = connectToDBP(*cfg)
}

func GetDB() *gorm.DB {
	dbOnce.Do(initDB)
	return db
}
