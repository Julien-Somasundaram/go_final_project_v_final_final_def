package server

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	cmd2 "github.com/Julien-Somasundaram/urlshortener/cmd"
	"github.com/Julien-Somasundaram/urlshortener/internal/api"
	"github.com/Julien-Somasundaram/urlshortener/internal/models"
	"github.com/Julien-Somasundaram/urlshortener/internal/monitor"
	"github.com/Julien-Somasundaram/urlshortener/internal/repository"
	"github.com/Julien-Somasundaram/urlshortener/internal/services"
	"github.com/Julien-Somasundaram/urlshortener/internal/workers"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var RunServerCmd = &cobra.Command{
	Use:   "run-server",
	Short: "Lance le serveur API de raccourcissement d'URLs et les processus de fond.",
	Long: `Cette commande initialise la base de données, configure les APIs,
démarre les workers asynchrones pour les clics et le moniteur d'URLs,
puis lance le serveur HTTP.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := cmd2.Cfg
		if cfg == nil {
			log.Fatalln("❌ Configuration non initialisée.")
		}

		// Connexion DB
		db, err := gorm.Open(sqlite.Open(cfg.Database.Name), &gorm.Config{})
		if err != nil {
			log.Fatalf("❌ Échec connexion DB : %v", err)
		}

		// Migration
		if err := db.AutoMigrate(&models.Link{}, &models.Click{}); err != nil {
			log.Fatalf("❌ Échec migration DB : %v", err)
		}

		// Repositories
		linkRepo := repository.NewGormLinkRepository(db)
		clickRepo := repository.NewGormClickRepository(db)
		log.Println("✅ Repositories initialisés.")

		// Services
		linkService := services.NewLinkService(linkRepo)
		log.Println("✅ Services métiers initialisés.")

		// Channel + Workers
		api.ClickEventsChannel = make(chan models.ClickEvent, cfg.Analytics.BufferSize)
		workers.StartClickWorkers(1, api.ClickEventsChannel, clickRepo)
		log.Printf("✅ Channel d'événements de clic initialisé avec un buffer de %d. %d worker(s) de clics démarré(s).",
			cfg.Analytics.BufferSize, 1)

		// Moniteur
		monitorInterval := time.Duration(cfg.Monitor.IntervalMinutes) * time.Minute
		urlMonitor := monitor.NewUrlMonitor(linkRepo, monitorInterval)
		go urlMonitor.Start()
		log.Printf("🛰️  Moniteur d'URLs démarré avec un intervalle de %v.", monitorInterval)

		// Routes
		router := gin.Default()
		api.SetupRoutes(router, linkService, cfg)
		log.Println("✅ Routes API configurées.")

		serverAddr := fmt.Sprintf(":%d", cfg.Server.Port)
		srv := &http.Server{
			Addr:    serverAddr,
			Handler: router,
		}

		// Serveur Gin dans une goroutine
		go func() {
			log.Printf("🚀 Serveur lancé sur %s ...", serverAddr)
			if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Fatalf("❌ Erreur serveur HTTP : %v", err)
			}
		}()

		// Shutdown propre
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Println("🛑 Signal d'arrêt reçu. Arrêt du serveur...")
		time.Sleep(5 * time.Second)
		log.Println("✅ Serveur arrêté proprement.")
	},
}

func init() {
	cmd2.RootCmd.AddCommand(RunServerCmd)
}
