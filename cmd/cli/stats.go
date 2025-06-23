package cli

import (
	"fmt"
	"log"
	"os"

	cmd2 "github.com/Julien-Somasundaram/urlshortener/cmd"
	"github.com/Julien-Somasundaram/urlshortener/internal/repository"
	"github.com/Julien-Somasundaram/urlshortener/internal/services"
	"github.com/spf13/cobra"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var shortCodeFlag string // Flag --code

var StatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Affiche les statistiques (nombre de clics) pour un lien court.",
	Long: `Cette commande permet de récupérer et d'afficher le nombre total de clics
pour une URL courte spécifique en utilisant son code.

Exemple:
  url-shortener stats --code="xyz123"`,
	Run: func(cmd *cobra.Command, args []string) {
		if shortCodeFlag == "" {
			fmt.Println("❌ Erreur : le flag --code est requis.")
			os.Exit(1)
		}

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
			log.Fatalf("❌ Échec récupération connexion DB : %v", err)
		}
		defer sqlDB.Close()

		linkRepo := repository.NewGormLinkRepository(db)
		// clickRepo := repository.NewGormClickRepository(db)
		linkService := services.NewLinkService(linkRepo)

		link, totalClicks, err := linkService.GetLinkStats(shortCodeFlag)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				fmt.Printf("❌ Aucun lien trouvé pour le code : %s\n", shortCodeFlag)
				os.Exit(1)
			}
			log.Fatalf("❌ Erreur lors de la récupération des stats : %v", err)
		}

		fmt.Printf("📊 Statistiques pour le code court : %s\n", link.ShortCode)
		fmt.Printf("🔗 URL longue : %s\n", link.LongURL)
		fmt.Printf("👁️  Total de clics : %d\n", totalClicks)
	},
}

func init() {
	StatsCmd.Flags().StringVar(&shortCodeFlag, "code", "", "Code court de l'URL à analyser")
	StatsCmd.MarkFlagRequired("code")
	cmd2.RootCmd.AddCommand(StatsCmd)
}
