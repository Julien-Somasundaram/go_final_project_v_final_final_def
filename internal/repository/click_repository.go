package repository

import (
	"github.com/Julien-Somasundaram/urlshortener/internal/models"
	"gorm.io/gorm"
)

// ClickRepository définit les opérations sur les clics.
type ClickRepository interface {
	CreateClick(click *models.Click) error
	CountClicksByLinkID(linkID uint) (int, error)
}

// GormClickRepository implémente ClickRepository avec GORM.
type GormClickRepository struct {
	db *gorm.DB
}

// NewClickRepository crée un nouveau dépôt GORM pour les clics.
func NewGormClickRepository(db *gorm.DB) *GormClickRepository {

	return &GormClickRepository{db: db}
}

// CreateClick insère un enregistrement de clic.
func (r *GormClickRepository) CreateClick(click *models.Click) error {
	result := r.db.Create(click)
	return result.Error
}

// CountClicksByLinkID retourne le nombre de clics pour un lien.
func (r *GormClickRepository) CountClicksByLinkID(linkID uint) (int, error) {
	var count int64
	result := r.db.Model(&models.Click{}).Where("link_id = ?", linkID).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return int(count), nil
}
