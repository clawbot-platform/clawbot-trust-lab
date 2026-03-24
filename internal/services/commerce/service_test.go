package commerce

import (
	"testing"
	"time"

	"clawbot-trust-lab/internal/domain/actors"
	domaincommerce "clawbot-trust-lab/internal/domain/commerce"
	"clawbot-trust-lab/internal/platform/store"
)

func TestCreateOrderComputesTotal(t *testing.T) {
	world := store.NewCommerceWorldStore()
	service := NewService(world)
	service.PutProduct(domaincommerce.Product{
		ID:         "product-1",
		MerchantID: "merchant-1",
		Name:       "Test Product",
		Amount:     4200,
		Currency:   "USD",
	})

	order, err := service.CreateOrder(CreateOrderInput{
		ID:                 "order-1",
		BuyerID:            "buyer-1",
		MerchantID:         "merchant-1",
		ProductIDs:         []string{"product-1"},
		SubmittedByActorID: "agent-1",
		DelegationMode:     actors.DelegationModeAgentAssisted,
		Status:             domaincommerce.OrderStatusAccepted,
		CreatedAt:          time.Date(2026, 3, 24, 9, 0, 0, 0, time.UTC),
		UpdatedAt:          time.Date(2026, 3, 24, 9, 1, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("CreateOrder() error = %v", err)
	}
	if order.TotalAmount != 4200 {
		t.Fatalf("expected total 4200, got %d", order.TotalAmount)
	}
}
