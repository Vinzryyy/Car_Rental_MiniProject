package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestXenditPaymentService_CreateInvoice(t *testing.T) {
	t.Run("mock path when secret key is empty", func(t *testing.T) {
		// Manually create service to ensure SecretKey is empty regardless of ENV
		service := &xenditPaymentService{
			config: &XenditConfig{
				SecretKey: "",
			},
		}
		
		url, err := service.CreateInvoice(context.Background(), "ORDER-1", 100.0, "test@example.com", "Test")
		
		assert.NoError(t, err)
		assert.Contains(t, url, "ORDER-1")
	})
}

func TestXenditPaymentService_VerifyPaymentNotification(t *testing.T) {
	service := &xenditPaymentService{
		config: &XenditConfig{
			CallbackToken: "valid-token",
		},
	}

	t.Run("valid token", func(t *testing.T) {
		isValid := service.VerifyPaymentNotification(context.Background(), "ORDER-1", "valid-token")
		assert.True(t, isValid)
	})

	t.Run("invalid token", func(t *testing.T) {
		isValid := service.VerifyPaymentNotification(context.Background(), "ORDER-1", "wrong-token")
		assert.False(t, isValid)
	})

	t.Run("empty token when required", func(t *testing.T) {
		isValid := service.VerifyPaymentNotification(context.Background(), "ORDER-1", "")
		assert.False(t, isValid)
	})

	t.Run("allow all when token not configured", func(t *testing.T) {
		serviceNoToken := &xenditPaymentService{
			config: &XenditConfig{
				CallbackToken: "",
			},
		}
		isValid := serviceNoToken.VerifyPaymentNotification(context.Background(), "ORDER-1", "any-token")
		assert.True(t, isValid)
	})
}
