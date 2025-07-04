package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/Julien-Somasundaram/urlshortener/internal/config"
	"github.com/spf13/cobra"
)

// Cfg est la variable globale qui contiendra la configuration chargée.
var Cfg *config.Config

// RootCmd représente la commande de base.
var RootCmd = &cobra.Command{
	Use:   "url-shortener",
	Short: "Un service de raccourcissement d'URLs avec API REST et CLI",
	Long: `url-shortener est une application complète pour gérer des URLs courtes.
Elle inclut un serveur API pour le raccourcissement et la redirection,
ainsi qu'une interface en ligne de commande pour l'administration.

Utilisez 'url-shortener [command] --help' pour plus d'informations sur une commande.`,
}

// Execute est le point d'entrée principal pour Cobra
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Erreur lors de l'exécution de la commande: %v\n", err)
		os.Exit(1)
	}
}

// Fonction spéciale exécutée avant le main()
func init() {
	// Initialisation automatique de la configuration au démarrage de Cobra
	cobra.OnInitialize(initConfig)

}

// initConfig charge la configuration avec Viper
func initConfig() {
	var err error
	Cfg, err = config.LoadConfig()
	if err != nil {
		log.Printf("⚠️  Attention: problème lors du chargement de la configuration: %v. Utilisation des valeurs par défaut.", err)
	}
}
