package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	//nolint:depguard
	"github.com/ElenaGrishkova/OtusGolangHomeWork/hw12_13_14_15_calendar/internal/app"
	//nolint:depguard
	"github.com/ElenaGrishkova/OtusGolangHomeWork/hw12_13_14_15_calendar/internal/logger"
	//nolint:depguard
	internalhttp "github.com/ElenaGrishkova/OtusGolangHomeWork/hw12_13_14_15_calendar/internal/server/http"
	//nolint:depguard
	memorystorage "github.com/ElenaGrishkova/OtusGolangHomeWork/hw12_13_14_15_calendar/internal/storage/memory"
	//nolint:depguard
	sqlstorage "github.com/ElenaGrishkova/OtusGolangHomeWork/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.json", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	// Создание конфига
	config := NewConfig()
	ctxConfig := context.TODO()
	if err := LoadConfig(ctxConfig, config, configFile); err != nil {
		log.Fatal(err)
	}

	// Логирование
	logg, err := logger.New(config.Logger.Level)
	if err != nil {
		log.Fatal(err)
	}

	// БД
	var storage app.Storage
	switch config.Database.Storage {
	case "in-memory":
		storage = memorystorage.NewMemoryStorage()
	case "database":
		driver := config.Database.Driver
		dsn := config.Database.Dsn
		storage = sqlstorage.NewSQLStorage(driver, dsn)
		if err := storage.Connect(); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatal("unknown storage type")
	}

	// Запуск приложения
	calendar := app.New(logg, storage)
	server := internalhttp.NewServer(net.JoinHostPort(config.Server.Host, config.Server.Port), logg, calendar)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
