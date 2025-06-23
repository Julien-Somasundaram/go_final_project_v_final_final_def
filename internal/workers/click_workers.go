package workers

import (
	"log"
	"time"

	"github.com/Julien-Somasundaram/urlshortener/internal/models"
	"github.com/Julien-Somasundaram/urlshortener/internal/repository"
)

// StartClickWorkers démarre plusieurs workers en parallèle
func StartClickWorkers(workerCount int, clickEventsChan <-chan models.ClickEvent, clickRepo repository.ClickRepository) {
	log.Printf("Starting %d click worker(s)...", workerCount)
	for i := 0; i < workerCount; i++ {
		go clickWorker(clickEventsChan, clickRepo)
	}
}

// Un worker écoute indéfiniment le channel et traite les événements
func clickWorker(clickEventsChan <-chan models.ClickEvent, clickRepo repository.ClickRepository) {
	for event := range clickEventsChan {
		click := &models.Click{
			LinkID:    event.LinkID,
			UserAgent: event.UserAgent,
			IPAddress: event.IPAddress,
			Timestamp: event.Timestamp,
		}

		err := clickRepo.CreateClick(click)
		if err != nil {
			log.Printf("ERROR: Failed to save click for LinkID %d (UserAgent: %s, IP: %s): %v",
				event.LinkID, event.UserAgent, event.IPAddress, err)
		} else {
			log.Printf("✅ Click recorded for LinkID %d at %v", event.LinkID, time.Now())
		}
	}
}
