package postgre

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	_ "github.com/lib/pq"
)

type Storage struct {
	DB *sql.DB
}

type PostgreConfig struct {
	Host     string `yaml:"host"`
	Port     int16  `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Dbname   string `yaml:"dbname"`
}

func FillConfig() PostgreConfig {
	cfg := PostgreConfig{}
	configPath := os.Getenv("CDPGPATH")
	if configPath == "" {
		log.Fatal("CDPGPATH is not set")
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file does not exist by path:%s", configPath)
	}
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config:%s", err)
	}
	return cfg
}
func SetUpConnect() Storage {
	cfg := FillConfig()
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Could not open connection to database:%s", err.Error())
	}
	stg := Storage{DB: db}
	CreateInitialTables(&stg)
	return stg

}

func CreateInitialTables(stg *Storage) {
	stmt, err := stg.DB.Prepare(`CREATE TABLE IF NOT EXISTS URL(
		id INTEGER PRIMARY KEY,
		alias TEXT NOT NULL UNIQUE,
		url TEXT NOT NULL);
		
	`)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := stmt.Exec(); err != nil {
		log.Fatal(err)
	}
	initIndex, err2 := stg.DB.Prepare(`CREATE INDEX IF NOT EXISTS idx_alias ON url(alias)`)
	if err2 != nil {
		log.Fatal(err2)
	}
	if _, err := initIndex.Exec(); err != nil {
		log.Fatal(err)
	}
}

func (stg *Storage) SaveUrl(urlToSave, alias string) error {

	insertSTMT, err := stg.DB.Prepare(`
	INSERT INTO url (url, alias) VALUES($1,$2)`)
	if err != nil {
		return errors.New("ошибка на инсерте")
	}
	if _, err = insertSTMT.Exec(urlToSave, alias); err != nil {
		return err
	}
	return nil
}
func (stg *Storage) GetUrl(alias string) (string, error) {
	SelectedRows, err := stg.DB.Query(`
	SELECT url FROM url WHERE alias = %s`, alias)
	if err != nil {
		return "", err
	}
	var url string
	for SelectedRows.Next() {
		if err := SelectedRows.Scan(&url); err != nil {
			return "", err
		}
	}
	return url, nil
}
