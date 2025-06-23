package api

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/Julien-Somasundaram/urlshortener/internal/config"
	"github.com/Julien-Somasundaram/urlshortener/internal/models"
	"github.com/Julien-Somasundaram/urlshortener/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var ClickEventsChannel chan models.ClickEvent // TODO 1: Channel global

// SetupRoutes configure toutes les routes de l'API
func SetupRoutes(router *gin.Engine, linkService *services.LinkService, cfg *config.Config) {
	if ClickEventsChannel == nil {
		ClickEventsChannel = make(chan models.ClickEvent, cfg.Analytics.BufferSize)
	}

	// Route de health check
	router.GET("/health", HealthCheckHandler)

	// Routes API REST
	api := router.Group("/api/v1")
	{
		api.POST("/links", CreateShortLinkHandler(linkService, cfg))
		api.GET("/links/:shortCode/stats", GetLinkStatsHandler(linkService))
	}

	// Redirection
	router.GET("/:shortCode", RedirectHandler(linkService, cfg))
}

// ───── HANDLERS ─────────────────────────────

func HealthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// Représente le corps d'une requête POST /links
type CreateLinkRequest struct {
	LongURL string `json:"long_url" binding:"required,url"`
}

func CreateShortLinkHandler(linkService *services.LinkService, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateLinkRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "URL invalide ou manquante"})
			return
		}

		link, err := linkService.CreateLink(req.LongURL)
		if err != nil {
			log.Printf("Erreur création lien: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur serveur"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"short_code":     link.ShortCode,
			"long_url":       link.LongURL,
			"full_short_url": cfg.Server.BaseURL + "/" + link.ShortCode,
		})
	}
}

func RedirectHandler(linkService *services.LinkService, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		shortCode := c.Param("shortCode")

		link, err := linkService.GetLinkByShortCode(shortCode)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Lien non trouvé"})
				return
			}
			log.Printf("Erreur récupération redirection: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur serveur"})
			return
		}

		clickEvent := models.ClickEvent{
			LinkID:    link.ID,
			Timestamp: time.Now(),
			UserAgent: c.Request.UserAgent(),
			IPAddress: c.ClientIP(),
		}

		// Multiplexage non bloquant
		select {
		case ClickEventsChannel <- clickEvent:
		default:
			log.Printf("⚠️  ClickEventsChannel is full, dropping click for %s", shortCode)
		}

		c.Redirect(http.StatusFound, link.LongURL)
	}
}

func GetLinkStatsHandler(linkService *services.LinkService) gin.HandlerFunc {
	return func(c *gin.Context) {
		shortCode := c.Param("shortCode")

		link, totalClicks, err := linkService.GetLinkStats(shortCode)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Lien non trouvé"})
				return
			}
			log.Printf("Erreur récupération stats: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur serveur"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"short_code":   link.ShortCode,
			"long_url":     link.LongURL,
			"total_clicks": totalClicks,
		})
	}
}
