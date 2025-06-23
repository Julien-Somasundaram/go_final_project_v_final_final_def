package cli

import (
	"fmt"
	"log"
	"net/url"
	"os"

	cmd2 "github.com/Julien-Somasundaram/urlshortener/cmd"
	"github.com/Julien-Somasundaram/urlshortener/internal/repository"
	"github.com/Julien-Somasundaram/urlshortener/internal/services"
	"github.com/spf13/cobra"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var longURLFlag string // --url

var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Crée une URL courte à partir d'une URL longue.",
	Long: `Cette commande raccourcit une URL longue fournie et affiche le code court généré.

Exemple:
  url-shortener create --url="https://www.google.com/search?q=go+lang"`,
	Run: func(cmd *cobra.Command, args []string) {
		if longURLFlag == "" {
			fmt.Println("❌ Le flag --url est requis.")
			os.Exit(1)
		}

		// Validation format URL
		if _, err := url.ParseRequestURI(longURLFlag); err != nil {
			fmt.Printf("❌ L'URL fournie est invalide : %v\n", err)
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
			log.Fatalf("❌ Échec récupération connexion SQL : %v", err)
		}
		defer sqlDB.Close()

		linkRepo := repository.NewGormLinkRepository(db)
		// clickRepo := repository.NewGormClickRepository(db)
		linkService := services.NewLinkService(linkRepo)

		link, err := linkService.CreateLink(longURLFlag)
		if err != nil {
			log.Fatalf("❌ Erreur lors de la création du lien court : %v", err)
		}

		fullShortURL := fmt.Sprintf("%s/%s", cfg.Server.BaseURL, link.ShortCode)
		fmt.Println("✅ URL courte créée avec succès :")
		fmt.Printf("🔗 Code : %s\n", link.ShortCode)
		fmt.Printf("🌐 URL complète : %s\n", fullShortURL)
	},
}

func init() {
	CreateCmd.Flags().StringVar(&longURLFlag, "url", "", "URL longue à raccourcir")
	CreateCmd.MarkFlagRequired("url")
	cmd2.RootCmd.AddCommand(CreateCmd)
}
