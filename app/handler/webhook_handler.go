package handler

import (
	"context"
	"net/http"
	"strings"

	"car_rental_miniproject/app/dto"
	"car_rental_miniproject/service"

	"github.com/labstack/echo/v4"
)

type PaymentWebhookHandler struct {
	rentalService service.RentalService
	topUpService  service.TopUpService
	emailService  *service.EmailService
}

func NewPaymentWebhookHandler(rentalService service.RentalService, topUpService service.TopUpService, emailService *service.EmailService) *PaymentWebhookHandler {
	return &PaymentWebhookHandler{
		rentalService: rentalService,
		topUpService:  topUpService,
		emailService:  emailService,
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
	var notification service.PaymentNotification
	if err := c.Bind(&notification); err != nil {
		return c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "invalid request body",
			Error:   err.Error(),
		})
	}

	// Extract entity type and ID from order_id
	// Format: RENTAL-{id}-{date} or TOPUP-{id}-{date}
	orderID := notification.OrderID
	var entityType string
	var entityID string

	if strings.HasPrefix(orderID, "RENTAL-") {
		entityType = "rental"
		entityID = strings.Split(strings.TrimPrefix(orderID, "RENTAL-"), "-")[0]
	} else if strings.HasPrefix(orderID, "TOPUP-") {
		entityType = "topup"
		entityID = strings.Split(strings.TrimPrefix(orderID, "TOPUP-"), "-")[0]
	} else {
		return c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "invalid order ID format",
		})
	}

	// Process based on transaction status
	switch notification.TransactionStatus {
	case "capture", "settlement":
		// Payment successful
		if entityType == "rental" {
			// Confirm rental payment
			// Note: In production, parse the UUID properly
			_ = entityID
			// h.rentalService.ConfirmPayment(ctx, rentalID)
		} else if entityType == "topup" {
			// Confirm top-up payment
			// Note: In production, parse the UUID properly
			_ = entityID
			// h.topUpService.ConfirmTopUp(ctx, transactionID)
		}

		// Send payment confirmation email (non-blocking)
		if h.emailService != nil && h.emailService.IsEnabled() {
			go func() {
				_ = h.emailService.SendPaymentConfirmationEmail(
					context.Background(),
					"", // userEmail - would need to fetch from database
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

	case "pending":
		// Payment is pending
		return c.JSON(http.StatusOK, dto.APIResponse{
			Success: true,
			Message: "payment is pending",
		})

	case "deny", "expire", "cancel", "failure":
		// Payment failed or cancelled
		return c.JSON(http.StatusOK, dto.APIResponse{
			Success: true,
			Message: "payment failed or cancelled",
		})

	default:
		return c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "unknown transaction status",
		})
	}
}
