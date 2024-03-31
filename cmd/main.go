package main

import (
	config "cdProjectGo/internal"
	"cdProjectGo/internal/storage/postgre"
	"fmt"
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
)

func main() {

	cfg := config.FillConfig()
	log := startLogger(cfg.Env)
	log.Debug("Запустили конфиг и логер")
	database := postgre.SetUpConnect()
	log.Debug("Подняли подклюение к базе")
	defer database.DB.Close()
	err := database.SaveUrl("https://vk.com/defindd", "вк")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(database.GetUrl("вк"))
	err = database.SaveUrl("https://postgrespro.ru/docs/postgresql/9.6/locale", "вк")
	if err != nil {
		fmt.Println(err)
		return
	}
}
func startLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log
}
