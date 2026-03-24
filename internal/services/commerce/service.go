package commerce

import (
	"fmt"
	"time"

	"clawbot-trust-lab/internal/domain/actors"
	"clawbot-trust-lab/internal/domain/commerce"
)

type Store interface {
	PutBuyer(commerce.Buyer)
	PutMerchant(commerce.Merchant)
	PutProduct(commerce.Product)
	GetProduct(string) (commerce.Product, error)
	PutOrder(commerce.Order)
	PutPayment(commerce.Payment)
	PutRefund(commerce.Refund)
	ListOrders() []commerce.Order
	GetOrder(string) (commerce.Order, error)
}

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{store: store}
}

func (s *Service) PutBuyer(item commerce.Buyer)       { s.store.PutBuyer(item) }
func (s *Service) PutMerchant(item commerce.Merchant) { s.store.PutMerchant(item) }
func (s *Service) PutProduct(item commerce.Product)   { s.store.PutProduct(item) }

type CreateOrderInput struct {
	ID                 string
	BuyerID            string
	MerchantID         string
	ProductIDs         []string
	SubmittedByActorID string
	DelegationMode     actors.DelegationMode
	MandateRef         string
	ProvenanceRef      string
	Status             commerce.OrderStatus
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

func (s *Service) CreateOrder(input CreateOrderInput) (commerce.Order, error) {
	if input.ID == "" {
		return commerce.Order{}, fmt.Errorf("order id is required")
	}
	if input.BuyerID == "" || input.MerchantID == "" {
		return commerce.Order{}, fmt.Errorf("buyer_id and merchant_id are required")
	}
	if len(input.ProductIDs) == 0 {
		return commerce.Order{}, fmt.Errorf("product_ids must not be empty")
	}

	var total int64
	currency := ""
	for _, productID := range input.ProductIDs {
		product, err := s.store.GetProduct(productID)
		if err != nil {
			return commerce.Order{}, err
		}
		total += product.Amount
		if currency == "" {
			currency = product.Currency
		}
	}

	item := commerce.Order{
		ID:                 input.ID,
		BuyerID:            input.BuyerID,
		MerchantID:         input.MerchantID,
		ProductIDs:         append([]string(nil), input.ProductIDs...),
		TotalAmount:        total,
		Currency:           currency,
		Status:             input.Status,
		SubmittedByActorID: input.SubmittedByActorID,
		DelegationMode:     input.DelegationMode,
		MandateRef:         input.MandateRef,
		ProvenanceRef:      input.ProvenanceRef,
		CreatedAt:          input.CreatedAt,
		UpdatedAt:          input.UpdatedAt,
	}
	s.store.PutOrder(item)
	return item, nil
}

type CreatePaymentInput struct {
	ID           string
	OrderID      string
	Amount       int64
	Currency     string
	Status       commerce.PaymentStatus
	Method       string
	AuthorizedAt time.Time
}

func (s *Service) CreatePayment(input CreatePaymentInput) commerce.Payment {
	item := commerce.Payment{
		ID:           input.ID,
		OrderID:      input.OrderID,
		Amount:       input.Amount,
		Currency:     input.Currency,
		Status:       input.Status,
		Method:       input.Method,
		AuthorizedAt: input.AuthorizedAt,
	}
	s.store.PutPayment(item)
	return item
}

type CreateRefundInput struct {
	ID                 string
	OrderID            string
	Amount             int64
	Status             commerce.RefundStatus
	RequestedByActorID string
	Reason             string
	CreatedAt          time.Time
}

func (s *Service) CreateRefund(input CreateRefundInput) commerce.Refund {
	item := commerce.Refund{
		ID:                 input.ID,
		OrderID:            input.OrderID,
		Amount:             input.Amount,
		Status:             input.Status,
		RequestedByActorID: input.RequestedByActorID,
		Reason:             input.Reason,
		CreatedAt:          input.CreatedAt,
	}
	s.store.PutRefund(item)
	return item
}

func (s *Service) UpdateOrder(item commerce.Order) {
	s.store.PutOrder(item)
}

func (s *Service) ListOrders() []commerce.Order {
	return s.store.ListOrders()
}

func (s *Service) GetOrder(id string) (commerce.Order, error) {
	return s.store.GetOrder(id)
}
