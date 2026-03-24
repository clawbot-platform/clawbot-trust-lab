package commerce

import (
	"time"

	"clawbot-trust-lab/internal/domain/actors"
)

type Buyer struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	RiskTier string   `json:"risk_tier"`
	Tags     []string `json:"tags"`
}

type Merchant struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Category string   `json:"category"`
	Tags     []string `json:"tags"`
}

type Product struct {
	ID         string   `json:"id"`
	MerchantID string   `json:"merchant_id"`
	Name       string   `json:"name"`
	Amount     int64    `json:"amount"`
	Currency   string   `json:"currency"`
	Category   string   `json:"category"`
	Tags       []string `json:"tags"`
}

type OrderStatus string
type PaymentStatus string
type RefundStatus string

const (
	OrderStatusAccepted       OrderStatus = "accepted"
	OrderStatusRefundReview   OrderStatus = "refund_review"
	OrderStatusRefundRejected OrderStatus = "refund_rejected"
)

const (
	PaymentStatusAuthorized PaymentStatus = "authorized"
)

const (
	RefundStatusRequested RefundStatus = "requested"
	RefundStatusRejected  RefundStatus = "rejected"
)

type Order struct {
	ID                 string                `json:"id"`
	BuyerID            string                `json:"buyer_id"`
	MerchantID         string                `json:"merchant_id"`
	ProductIDs         []string              `json:"product_ids"`
	TotalAmount        int64                 `json:"total_amount"`
	Currency           string                `json:"currency"`
	Status             OrderStatus           `json:"status"`
	SubmittedByActorID string                `json:"submitted_by_actor_id"`
	DelegationMode     actors.DelegationMode `json:"delegation_mode"`
	MandateRef         string                `json:"mandate_ref"`
	ProvenanceRef      string                `json:"provenance_ref"`
	CreatedAt          time.Time             `json:"created_at"`
	UpdatedAt          time.Time             `json:"updated_at"`
}

type Payment struct {
	ID           string        `json:"id"`
	OrderID      string        `json:"order_id"`
	Amount       int64         `json:"amount"`
	Currency     string        `json:"currency"`
	Status       PaymentStatus `json:"status"`
	Method       string        `json:"method"`
	AuthorizedAt time.Time     `json:"authorized_at"`
}

type Refund struct {
	ID                 string       `json:"id"`
	OrderID            string       `json:"order_id"`
	Amount             int64        `json:"amount"`
	Status             RefundStatus `json:"status"`
	RequestedByActorID string       `json:"requested_by_actor_id"`
	Reason             string       `json:"reason"`
	CreatedAt          time.Time    `json:"created_at"`
}
