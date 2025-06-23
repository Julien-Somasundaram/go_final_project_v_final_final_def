package main

import (
	"github.com/Julien-Somasundaram/urlshortener/cmd"
	_ "github.com/Julien-Somasundaram/urlshortener/cmd/cli"    // Exécute les init() des sous-commandes CLI
	_ "github.com/Julien-Somasundaram/urlshortener/cmd/server" // Exécute les init() de la commande serveur
)

func main() {
	cmd.Execute()
}
