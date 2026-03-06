package service

import (
	"testing"

	"car_rental_miniproject/app/config"

	"github.com/stretchr/testify/assert"
)

func TestEmailService_IsEnabled(t *testing.T) {
	t.Run("enabled in config but no credentials", func(t *testing.T) {
		cfg := &config.Config{
			Email: config.EmailConfig{
				IsEnabled: true,
			},
		}
		service := NewEmailService(cfg)
		// Should be false because no credentials provided
		assert.False(t, service.IsEnabled())
	})

	t.Run("disabled in config", func(t *testing.T) {
		cfg := &config.Config{
			Email: config.EmailConfig{
				IsEnabled: false,
			},
		}
		service := NewEmailService(cfg)
		assert.False(t, service.IsEnabled())
	})
}

func TestEmailService_GetFromEmail(t *testing.T) {
	cfg := &config.Config{
		Email: config.EmailConfig{
			FromEmail: "noreply@example.com",
		},
	}
	service := NewEmailService(cfg)
	assert.Equal(t, "noreply@example.com", service.GetFromEmail())
}
