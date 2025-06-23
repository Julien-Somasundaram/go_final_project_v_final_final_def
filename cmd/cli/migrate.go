package cli

import (
	"fmt"
	"log"

	cmd2 "github.com/Julien-Somasundaram/urlshortener/cmd"
	"github.com/Julien-Somasundaram/urlshortener/internal/models"
	"github.com/spf13/cobra"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var MigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Exécute les migrations de la base de données pour créer ou mettre à jour les tables.",
	Long: `Cette commande se connecte à la base de données configurée (SQLite)
et exécute les migrations automatiques de GORM pour créer les tables 'links' et 'clicks'
basées sur les modèles Go.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := cmd2.Cfg
		if cfg == nil {
			log.Fatalln("❌ Configuration non initialisée.")
		}

		db, err := gorm.Open(sqlite.Open(cfg.Database.Name), &gorm.Config{})
		if err != nil {
			log.Fatalf("❌ Échec connexion DB : %v", err)
		}

		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("❌ Échec récupération connexion SQL : %v", err)
		}
		defer sqlDB.Close()

		if err := db.AutoMigrate(&models.Link{}, &models.Click{}); err != nil {
			log.Fatalf("❌ Erreur migration : %v", err)
		}

		fmt.Println("✅ Migrations de la base de données exécutées avec succès.")
	},
}

func init() {
	cmd2.RootCmd.AddCommand(MigrateCmd)
}
