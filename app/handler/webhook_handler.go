package handler

import (
	"context"
	"log"
	"net/http"
	"strings"

	"car_rental_miniproject/app/dto"
	"car_rental_miniproject/service"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type PaymentWebhookHandler struct {
	rentalService  service.RentalService
	topUpService   service.TopUpService
	paymentService service.PaymentService
	emailService   *service.EmailService
}

func NewPaymentWebhookHandler(rentalService service.RentalService, topUpService service.TopUpService, paymentService service.PaymentService, emailService *service.EmailService) *PaymentWebhookHandler {
	return &PaymentWebhookHandler{
		rentalService:  rentalService,
		topUpService:   topUpService,
		paymentService: paymentService,
		emailService:   emailService,
	}
}

// PaymentNotification godoc
// @Summary Payment notification webhook
// @Description Receive payment notifications from Xendit
// @Tags webhook
// @Accept json
// @Produce json
// @Param request body object true "Payment notification from Xendit"
// @Success 200 {object} dto.APIResponse
// @Failure 400 {object} dto.APIResponse
// @Router /api/webhook/payment [post]
func (h *PaymentWebhookHandler) PaymentNotification(c echo.Context) error {
	callbackToken := c.Request().Header.Get("x-callback-token")
	
	var notification service.PaymentNotification
	if err := c.Bind(&notification); err != nil {
		return c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "invalid request body",
			Error:   err.Error(),
		})
	}

	// Verify callback token
	if !h.paymentService.VerifyPaymentNotification(c.Request().Context(), notification.OrderID, callbackToken) {
		return c.JSON(http.StatusUnauthorized, dto.APIResponse{
			Success: false,
			Message: "unauthorized callback token",
		})
	}

	// Extract entity type and ID from order_id
	// Format: RENTAL-{uuid} or TOPUP-{uuid}
	orderID := notification.OrderID
	var entityType string
	var entityIDStr string

	if strings.HasPrefix(orderID, "RENTAL-") {
		entityType = "rental"
		entityIDStr = strings.TrimPrefix(orderID, "RENTAL-")
	} else if strings.HasPrefix(orderID, "TOPUP-") {
		entityType = "topup"
		entityIDStr = strings.TrimPrefix(orderID, "TOPUP-")
	} else {
		return c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "invalid order ID format",
		})
	}

	entityID, err := uuid.Parse(entityIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "invalid ID in order ID",
			Error:   err.Error(),
		})
	}

	// Process based on transaction status
	switch notification.TransactionStatus {
	case "PAID", "SETTLED", "capture", "settlement":
		// Payment successful
		if entityType == "rental" {
			// Confirm rental payment via external gateway
			err = h.rentalService.ConfirmExternalPayment(c.Request().Context(), entityID)
			if err != nil {
				log.Printf("Webhook error: failed to confirm rental %s: %v", entityID, err)
			}
		} else if entityType == "topup" {
			// Confirm top-up payment
			err = h.topUpService.ConfirmTopUp(c.Request().Context(), entityID)
			if err != nil {
				log.Printf("Webhook error: failed to confirm top-up %s: %v", entityID, err)
			}
		}

		// Send payment confirmation email (non-blocking)
		if h.emailService != nil && h.emailService.IsEnabled() {
			go func() {
				// In a real app, we'd fetch the user email here
				_ = h.emailService.SendPaymentConfirmationEmail(
					context.Background(),
					"customer@example.com", 
					"Customer",
					notification.OrderID,
					notification.GrossAmount,
					notification.TransactionStatus,
				)
			}()
		}

		return c.JSON(http.StatusOK, dto.APIResponse{
			Success: true,
			Message: "payment notification processed successfully",
		})

	case "PENDING":
		return c.JSON(http.StatusOK, dto.APIResponse{
			Success: true,
			Message: "payment is pending",
		})

	case "EXPIRED", "FAILED", "CANCELLED":
		// Handle failure if needed (e.g., mark rental as cancelled)
		return c.JSON(http.StatusOK, dto.APIResponse{
			Success: true,
			Message: "payment failed or expired",
		})

	default:
		return c.JSON(http.StatusOK, dto.APIResponse{
			Success: true,
			Message: "received status: " + notification.TransactionStatus,
		})
	}
}
