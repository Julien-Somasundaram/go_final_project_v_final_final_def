package services

import (
	"fmt"

	"github.com/Julien-Somasundaram/urlshortener/internal/models"
	"github.com/Julien-Somasundaram/urlshortener/internal/repository"
)

// ClickService fournit des méthodes métier pour les clics.
type ClickService struct {
	clickRepo repository.ClickRepository
}

// NewClickService crée un nouveau service de clics.
func NewClickService(clickRepo repository.ClickRepository) *ClickService {
	return &ClickService{
		clickRepo: clickRepo,
	}
}

// RecordClick enregistre un clic.
func (s *ClickService) RecordClick(click *models.Click) error {
	if err := s.clickRepo.CreateClick(click); err != nil {
		return fmt.Errorf("échec de l'enregistrement du clic : %w", err)
	}
	return nil
}

// GetClicksCountByLinkID retourne le nombre total de clics pour un lien.
func (s *ClickService) GetClicksCountByLinkID(linkID uint) (int, error) {
	return s.clickRepo.CountClicksByLinkID(linkID)
}
