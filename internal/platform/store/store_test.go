package store

import (
	"testing"
	"time"

	"clawbot-trust-lab/internal/domain/benchmark"
	"clawbot-trust-lab/internal/domain/commerce"
	"clawbot-trust-lab/internal/domain/detection"
	"clawbot-trust-lab/internal/domain/events"
	"clawbot-trust-lab/internal/domain/replay"
	"clawbot-trust-lab/internal/domain/trust"
)

func TestOperatorStoreReviewLifecycle(t *testing.T) {
	store := NewOperatorStore()
	review := benchmark.PromotionReview{PromotionID: "promo-1"}

	store.PutReview(review)

	got, ok := store.GetReview("promo-1")
	if !ok || got.PromotionID != "promo-1" {
		t.Fatalf("expected stored review, got %#v ok=%t", got, ok)
	}
	if len(store.ListReviews()) != 1 {
		t.Fatalf("expected 1 review, got %#v", store.ListReviews())
	}
}

func TestDetectionStorePreservesOrderAndSummary(t *testing.T) {
	store := NewDetectionStore()
	first := detection.DetectionResult{ID: "det-1", Status: detection.DetectionStatusClean}
	second := detection.DetectionResult{ID: "det-2", Status: detection.DetectionStatusBlocked}

	store.Put(first)
	store.Put(second)
	store.Put(first)

	items := store.List()
	if len(items) != 2 || items[0].ID != "det-1" || items[1].ID != "det-2" {
		t.Fatalf("unexpected detection order %#v", items)
	}
	got, err := store.Get("det-2")
	if err != nil || got.ID != "det-2" {
		t.Fatalf("expected det-2, got %#v err=%v", got, err)
	}
	if _, err := store.Get("missing"); err == nil {
		t.Fatal("expected missing detection lookup to fail")
	}
	summary := store.Summary()
	if summary.Total != 2 || summary.LastResultID != "det-1" || summary.TotalByStatus[detection.DetectionStatusBlocked] != 1 {
		t.Fatalf("unexpected detection summary %#v", summary)
	}
}

func TestTrustArtifactStoreListSorted(t *testing.T) {
	store := NewInMemoryTrustArtifactStore()
	if err := store.Create(trust.TrustArtifact{ID: "ta-2"}); err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if err := store.Create(trust.TrustArtifact{ID: "ta-1"}); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	items := store.List()
	if len(items) != 2 || items[0].ID != "ta-1" || items[1].ID != "ta-2" {
		t.Fatalf("unexpected trust artifact order %#v", items)
	}
}

func TestFileReplayStoreCreateListAndReload(t *testing.T) {
	dir := t.TempDir()
	store, err := NewFileReplayStore(dir)
	if err != nil {
		t.Fatalf("NewFileReplayStore() error = %v", err)
	}

	first := replay.ReplayCase{ID: "rc-1", RecordedAt: time.Date(2026, 3, 26, 10, 0, 0, 0, time.UTC)}
	second := replay.ReplayCase{ID: "rc-2", RecordedAt: time.Date(2026, 3, 26, 11, 0, 0, 0, time.UTC)}
	if _, err := store.Create(second); err != nil {
		t.Fatalf("Create(second) error = %v", err)
	}
	if _, err := store.Create(first); err != nil {
		t.Fatalf("Create(first) error = %v", err)
	}

	items := store.List()
	if len(items) != 2 || items[0].ID != "rc-1" || items[1].ID != "rc-2" {
		t.Fatalf("unexpected replay ordering %#v", items)
	}

	reloaded, err := NewFileReplayStore(dir)
	if err != nil {
		t.Fatalf("NewFileReplayStore(reload) error = %v", err)
	}
	reloadedItems := reloaded.List()
	if len(reloadedItems) != 2 || reloadedItems[0].ID != "rc-1" {
		t.Fatalf("unexpected reloaded replay items %#v", reloadedItems)
	}
}

func TestReplayFileNameValidation(t *testing.T) {
	for _, id := range []string{"", "../escape", "/abs/path", "nested/id"} {
		if _, err := replayFileName(id); err == nil {
			t.Fatalf("expected replayFileName(%q) to fail", id)
		}
	}
	if got, err := replayFileName("rc-1"); err != nil || got != "rc-1.json" {
		t.Fatalf("unexpected replay file name %q err=%v", got, err)
	}
}

func TestCommerceWorldStoreAccessors(t *testing.T) {
	store := NewCommerceWorldStore()
	store.PutBuyer(commerce.Buyer{ID: "buyer-1"})
	store.PutMerchant(commerce.Merchant{ID: "merchant-1"})
	store.PutProduct(commerce.Product{ID: "product-1"})
	store.PutOrder(commerce.Order{ID: "order-1"})
	store.PutPayment(commerce.Payment{ID: "payment-1"})
	store.PutRefund(commerce.Refund{ID: "refund-1"})
	store.PutMandate(trust.Mandate{ID: "mandate-1"})
	store.PutProvenance(trust.ProvenanceRecord{ID: "prov-1"})
	store.PutApproval(trust.ApprovalRecord{ID: "approval-1"})
	store.PutDecision(trust.TrustDecision{ID: "decision-1"})
	store.AppendEvent(events.Record{ID: "event-1"})

	if _, err := store.GetProduct("product-1"); err != nil {
		t.Fatalf("GetProduct() error = %v", err)
	}
	if _, err := store.GetOrder("order-1"); err != nil {
		t.Fatalf("GetOrder() error = %v", err)
	}
	if _, err := store.GetRefund("refund-1"); err != nil {
		t.Fatalf("GetRefund() error = %v", err)
	}
	if _, err := store.GetMandate("mandate-1"); err != nil {
		t.Fatalf("GetMandate() error = %v", err)
	}
	if _, err := store.GetProvenance("prov-1"); err != nil {
		t.Fatalf("GetProvenance() error = %v", err)
	}
	if _, err := store.GetTrustDecision("decision-1"); err != nil {
		t.Fatalf("GetTrustDecision() error = %v", err)
	}
	if len(store.ListRefunds()) != 1 || len(store.ListOrders()) != 1 || len(store.ListApprovals()) != 1 || len(store.ListTrustDecisions()) != 1 || len(store.ListEvents()) != 1 {
		t.Fatalf("unexpected commerce world listings")
	}
	if _, err := store.GetOrder("missing"); err == nil {
		t.Fatal("expected missing order lookup to fail")
	}
	if _, err := store.GetRefund("missing"); err == nil {
		t.Fatal("expected missing refund lookup to fail")
	}
	if _, err := store.GetProduct("missing"); err == nil {
		t.Fatal("expected missing product lookup to fail")
	}
}
