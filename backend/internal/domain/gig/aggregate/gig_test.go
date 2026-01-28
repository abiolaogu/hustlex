package aggregate

import (
	"testing"
	"time"

	"hustlex/internal/domain/shared/valueobject"
)

func TestNewBudget(t *testing.T) {
	min := valueobject.MustNewMoney(10000, valueobject.NGN)
	max := valueobject.MustNewMoney(50000, valueobject.NGN)

	budget, err := NewBudget(min, max)
	if err != nil {
		t.Fatalf("NewBudget() unexpected error: %v", err)
	}

	if budget.Min().Amount() != 10000 {
		t.Errorf("Budget.Min() = %d, want 10000", budget.Min().Amount())
	}
	if budget.Max().Amount() != 50000 {
		t.Errorf("Budget.Max() = %d, want 50000", budget.Max().Amount())
	}
}

func TestNewBudget_InvalidRange(t *testing.T) {
	min := valueobject.MustNewMoney(50000, valueobject.NGN)
	max := valueobject.MustNewMoney(10000, valueobject.NGN)

	_, err := NewBudget(min, max)
	if err != ErrInvalidBudget {
		t.Errorf("NewBudget() error = %v, want ErrInvalidBudget", err)
	}
}

func TestBudget_Contains(t *testing.T) {
	min := valueobject.MustNewMoney(10000, valueobject.NGN)
	max := valueobject.MustNewMoney(50000, valueobject.NGN)
	budget, _ := NewBudget(min, max)

	tests := []struct {
		amount int64
		want   bool
	}{
		{10000, true},  // Min
		{30000, true},  // Middle
		{50000, true},  // Max
		{5000, false},  // Below min
		{60000, false}, // Above max
	}

	for _, tt := range tests {
		amount := valueobject.MustNewMoney(tt.amount, valueobject.NGN)
		got := budget.Contains(amount)
		if got != tt.want {
			t.Errorf("Budget.Contains(%d) = %v, want %v", tt.amount, got, tt.want)
		}
	}
}

func TestNewProposal(t *testing.T) {
	proposalID := valueobject.GenerateProposalID()
	hustlerID := valueobject.GenerateUserID()
	price := valueobject.MustNewMoney(25000, valueobject.NGN)

	proposal := NewProposal(
		proposalID,
		hustlerID,
		"I am the best for this job",
		price,
		7,
		[]string{"portfolio.pdf"},
	)

	if !proposal.ID().Equals(proposalID) {
		t.Error("NewProposal() should set ID")
	}
	if !proposal.HustlerID().Equals(hustlerID) {
		t.Error("NewProposal() should set HustlerID")
	}
	if proposal.CoverLetter() != "I am the best for this job" {
		t.Errorf("NewProposal() cover letter = %s", proposal.CoverLetter())
	}
	if proposal.ProposedPrice().Amount() != 25000 {
		t.Errorf("NewProposal() price = %d, want 25000", proposal.ProposedPrice().Amount())
	}
	if proposal.DeliveryDays() != 7 {
		t.Errorf("NewProposal() delivery = %d, want 7", proposal.DeliveryDays())
	}
	if proposal.Status() != ProposalStatusPending {
		t.Errorf("NewProposal() status = %s, want pending", proposal.Status())
	}
	if !proposal.IsPending() {
		t.Error("NewProposal() IsPending should be true")
	}
}

func TestProposal_Accept(t *testing.T) {
	proposal := createTestProposal()

	proposal.Accept()

	if proposal.Status() != ProposalStatusAccepted {
		t.Errorf("Accept() status = %s, want accepted", proposal.Status())
	}
	if proposal.IsPending() {
		t.Error("Accept() IsPending should be false")
	}
}

func TestProposal_Reject(t *testing.T) {
	proposal := createTestProposal()

	proposal.Reject()

	if proposal.Status() != ProposalStatusRejected {
		t.Errorf("Reject() status = %s, want rejected", proposal.Status())
	}
}

func TestProposal_Withdraw(t *testing.T) {
	proposal := createTestProposal()

	proposal.Withdraw()

	if proposal.Status() != ProposalStatusWithdrawn {
		t.Errorf("Withdraw() status = %s, want withdrawn", proposal.Status())
	}
}

func TestNewGig(t *testing.T) {
	gigID := valueobject.GenerateGigID()
	clientID := valueobject.GenerateUserID()
	budget := createTestBudget()

	gig, err := NewGig(
		gigID,
		clientID,
		"Build a website",
		"I need a responsive website",
		"Web Development",
		budget,
		14,
		true,
	)

	if err != nil {
		t.Fatalf("NewGig() unexpected error: %v", err)
	}

	if !gig.ID().Equals(gigID) {
		t.Error("NewGig() should set ID")
	}
	if !gig.ClientID().Equals(clientID) {
		t.Error("NewGig() should set ClientID")
	}
	if gig.Title() != "Build a website" {
		t.Errorf("NewGig() title = %s", gig.Title())
	}
	if gig.Status() != GigStatusOpen {
		t.Errorf("NewGig() status = %s, want open", gig.Status())
	}
	if !gig.IsOpen() {
		t.Error("NewGig() IsOpen should be true")
	}
	if gig.IsRemote() != true {
		t.Error("NewGig() IsRemote should be true")
	}
	if gig.ViewCount() != 0 {
		t.Errorf("NewGig() viewCount = %d, want 0", gig.ViewCount())
	}
	if gig.ProposalCount() != 0 {
		t.Errorf("NewGig() proposalCount = %d, want 0", gig.ProposalCount())
	}

	events := gig.DomainEvents()
	if len(events) != 1 {
		t.Fatalf("NewGig() should record 1 event, got %d", len(events))
	}
}

func TestGig_Update(t *testing.T) {
	gig := createTestGig()
	gig.ClearEvents()
	newBudget := createTestBudget()
	deadline := time.Now().Add(30 * 24 * time.Hour)

	err := gig.Update(
		"Updated Title",
		"Updated Description",
		"Design",
		nil,
		newBudget,
		21,
		&deadline,
		false,
		"Lagos",
		[]string{"new.pdf"},
		[]string{"design", "ui"},
	)

	if err != nil {
		t.Fatalf("Update() unexpected error: %v", err)
	}

	if gig.Title() != "Updated Title" {
		t.Errorf("Update() title = %s, want Updated Title", gig.Title())
	}
	if gig.Description() != "Updated Description" {
		t.Errorf("Update() description = %s", gig.Description())
	}
	if gig.Category() != "Design" {
		t.Errorf("Update() category = %s, want Design", gig.Category())
	}
	if gig.Location() != "Lagos" {
		t.Errorf("Update() location = %s, want Lagos", gig.Location())
	}
}

func TestGig_Update_NotOpen(t *testing.T) {
	gig := createTestGig()
	gig.Cancel("test")

	err := gig.Update("Title", "Desc", "Cat", nil, createTestBudget(), 7, nil, true, "", nil, nil)
	if err != ErrCannotUpdateGig {
		t.Errorf("Update() on cancelled gig error = %v, want ErrCannotUpdateGig", err)
	}
}

func TestGig_IncrementViewCount(t *testing.T) {
	gig := createTestGig()

	gig.IncrementViewCount()
	gig.IncrementViewCount()
	gig.IncrementViewCount()

	if gig.ViewCount() != 3 {
		t.Errorf("IncrementViewCount() count = %d, want 3", gig.ViewCount())
	}
}

func TestGig_SetFeatured(t *testing.T) {
	gig := createTestGig()

	if gig.IsFeatured() {
		t.Error("new gig should not be featured")
	}

	gig.SetFeatured(true)

	if !gig.IsFeatured() {
		t.Error("SetFeatured(true) should make gig featured")
	}
}

func TestGig_Cancel(t *testing.T) {
	gig := createTestGig()
	gig.ClearEvents()

	err := gig.Cancel("Client changed mind")
	if err != nil {
		t.Fatalf("Cancel() unexpected error: %v", err)
	}

	if gig.Status() != GigStatusCancelled {
		t.Errorf("Cancel() status = %s, want cancelled", gig.Status())
	}
	if gig.IsOpen() {
		t.Error("Cancel() IsOpen should be false")
	}

	events := gig.DomainEvents()
	if len(events) != 1 {
		t.Fatalf("Cancel() should record 1 event, got %d", len(events))
	}
}

func TestGig_Cancel_NotOpen(t *testing.T) {
	gig := createTestGig()
	gig.MarkCompleted()

	err := gig.Cancel("reason")
	if err != ErrGigAlreadyInProgress {
		t.Errorf("Cancel() on completed gig error = %v, want ErrGigAlreadyInProgress", err)
	}
}

func TestGig_SubmitProposal(t *testing.T) {
	gig := createTestGig()
	gig.ClearEvents()

	proposal := createTestProposalWithPrice(25000)

	err := gig.SubmitProposal(proposal)
	if err != nil {
		t.Fatalf("SubmitProposal() unexpected error: %v", err)
	}

	if gig.ProposalCount() != 1 {
		t.Errorf("SubmitProposal() count = %d, want 1", gig.ProposalCount())
	}

	events := gig.DomainEvents()
	if len(events) != 1 {
		t.Fatalf("SubmitProposal() should record 1 event, got %d", len(events))
	}
}

func TestGig_SubmitProposal_NotOpen(t *testing.T) {
	gig := createTestGig()
	gig.Cancel("test")

	proposal := createTestProposalWithPrice(25000)

	err := gig.SubmitProposal(proposal)
	if err != ErrGigNotOpen {
		t.Errorf("SubmitProposal() on cancelled gig error = %v, want ErrGigNotOpen", err)
	}
}

func TestGig_SubmitProposal_SelfProposal(t *testing.T) {
	clientID := valueobject.GenerateUserID()
	gig := createTestGigWithClient(clientID)

	price := valueobject.MustNewMoney(25000, valueobject.NGN)
	proposal := NewProposal(
		valueobject.GenerateProposalID(),
		clientID, // Same as client
		"Self proposal",
		price,
		7,
		nil,
	)

	err := gig.SubmitProposal(proposal)
	if err != ErrCannotProposeSelf {
		t.Errorf("SubmitProposal() by client error = %v, want ErrCannotProposeSelf", err)
	}
}

func TestGig_SubmitProposal_Duplicate(t *testing.T) {
	gig := createTestGig()
	hustlerID := valueobject.GenerateUserID()

	price := valueobject.MustNewMoney(25000, valueobject.NGN)
	proposal1 := NewProposal(valueobject.GenerateProposalID(), hustlerID, "First", price, 7, nil)
	proposal2 := NewProposal(valueobject.GenerateProposalID(), hustlerID, "Second", price, 7, nil)

	gig.SubmitProposal(proposal1)
	err := gig.SubmitProposal(proposal2)

	if err != ErrAlreadyProposed {
		t.Errorf("SubmitProposal() duplicate error = %v, want ErrAlreadyProposed", err)
	}
}

func TestGig_SubmitProposal_PriceBelowBudget(t *testing.T) {
	gig := createTestGig() // Budget: 10000-50000

	price := valueobject.MustNewMoney(5000, valueobject.NGN) // Below min
	proposal := NewProposal(valueobject.GenerateProposalID(), valueobject.GenerateUserID(), "Test", price, 7, nil)

	err := gig.SubmitProposal(proposal)
	if err != ErrPriceBelowBudget {
		t.Errorf("SubmitProposal() below budget error = %v, want ErrPriceBelowBudget", err)
	}
}

func TestGig_SubmitProposal_PriceAboveBudget(t *testing.T) {
	gig := createTestGig() // Budget: 10000-50000

	price := valueobject.MustNewMoney(60000, valueobject.NGN) // Above max
	proposal := NewProposal(valueobject.GenerateProposalID(), valueobject.GenerateUserID(), "Test", price, 7, nil)

	err := gig.SubmitProposal(proposal)
	if err != ErrPriceAboveBudget {
		t.Errorf("SubmitProposal() above budget error = %v, want ErrPriceAboveBudget", err)
	}
}

func TestGig_WithdrawProposal(t *testing.T) {
	gig := createTestGig()
	hustlerID := valueobject.GenerateUserID()
	proposalID := valueobject.GenerateProposalID()
	price := valueobject.MustNewMoney(25000, valueobject.NGN)
	proposal := NewProposal(proposalID, hustlerID, "Test", price, 7, nil)
	gig.SubmitProposal(proposal)
	gig.ClearEvents()

	err := gig.WithdrawProposal(proposalID, hustlerID)
	if err != nil {
		t.Fatalf("WithdrawProposal() unexpected error: %v", err)
	}

	found := gig.FindProposal(proposalID)
	if found.Status() != ProposalStatusWithdrawn {
		t.Errorf("WithdrawProposal() status = %s, want withdrawn", found.Status())
	}
}

func TestGig_WithdrawProposal_NotFound(t *testing.T) {
	gig := createTestGig()

	err := gig.WithdrawProposal(valueobject.GenerateProposalID(), valueobject.GenerateUserID())
	if err != ErrProposalNotFound {
		t.Errorf("WithdrawProposal() not found error = %v, want ErrProposalNotFound", err)
	}
}

func TestGig_WithdrawProposal_NotOwner(t *testing.T) {
	gig := createTestGig()
	proposalID := valueobject.GenerateProposalID()
	price := valueobject.MustNewMoney(25000, valueobject.NGN)
	proposal := NewProposal(proposalID, valueobject.GenerateUserID(), "Test", price, 7, nil)
	gig.SubmitProposal(proposal)

	differentUser := valueobject.GenerateUserID()
	err := gig.WithdrawProposal(proposalID, differentUser)
	if err != ErrNotGigOwner {
		t.Errorf("WithdrawProposal() not owner error = %v, want ErrNotGigOwner", err)
	}
}

func TestGig_AcceptProposal(t *testing.T) {
	gig := createTestGig()
	proposalID := valueobject.GenerateProposalID()
	hustlerID := valueobject.GenerateUserID()
	price := valueobject.MustNewMoney(25000, valueobject.NGN)
	proposal := NewProposal(proposalID, hustlerID, "Test", price, 7, nil)
	gig.SubmitProposal(proposal)
	gig.ClearEvents()

	contractID := valueobject.GenerateContractID()
	data, err := gig.AcceptProposal(proposalID, contractID)

	if err != nil {
		t.Fatalf("AcceptProposal() unexpected error: %v", err)
	}

	if !data.ContractID.Equals(contractID) {
		t.Error("AcceptProposal() should return contract ID")
	}
	if !data.HustlerID.Equals(hustlerID) {
		t.Error("AcceptProposal() should return hustler ID")
	}
	if data.AgreedPrice.Amount() != 25000 {
		t.Errorf("AcceptProposal() price = %d, want 25000", data.AgreedPrice.Amount())
	}
	if data.PlatformFee != 2500 { // 10%
		t.Errorf("AcceptProposal() fee = %d, want 2500", data.PlatformFee)
	}
	if gig.Status() != GigStatusInProgress {
		t.Errorf("AcceptProposal() status = %s, want in_progress", gig.Status())
	}
	if gig.AcceptedProposalID() == nil || !gig.AcceptedProposalID().Equals(proposalID) {
		t.Error("AcceptProposal() should set accepted proposal ID")
	}

	// Verify proposal is accepted
	accepted := gig.FindProposal(proposalID)
	if accepted.Status() != ProposalStatusAccepted {
		t.Errorf("AcceptProposal() proposal status = %s, want accepted", accepted.Status())
	}
}

func TestGig_AcceptProposal_RejectsOthers(t *testing.T) {
	gig := createTestGig()

	// Submit 3 proposals
	proposalID1 := valueobject.GenerateProposalID()
	proposalID2 := valueobject.GenerateProposalID()
	proposalID3 := valueobject.GenerateProposalID()
	price := valueobject.MustNewMoney(25000, valueobject.NGN)

	gig.SubmitProposal(NewProposal(proposalID1, valueobject.GenerateUserID(), "P1", price, 7, nil))
	gig.SubmitProposal(NewProposal(proposalID2, valueobject.GenerateUserID(), "P2", price, 7, nil))
	gig.SubmitProposal(NewProposal(proposalID3, valueobject.GenerateUserID(), "P3", price, 7, nil))

	// Accept proposal 2
	gig.AcceptProposal(proposalID2, valueobject.GenerateContractID())

	// Check statuses
	p1 := gig.FindProposal(proposalID1)
	p2 := gig.FindProposal(proposalID2)
	p3 := gig.FindProposal(proposalID3)

	if p1.Status() != ProposalStatusRejected {
		t.Errorf("Other proposal 1 status = %s, want rejected", p1.Status())
	}
	if p2.Status() != ProposalStatusAccepted {
		t.Errorf("Accepted proposal status = %s, want accepted", p2.Status())
	}
	if p3.Status() != ProposalStatusRejected {
		t.Errorf("Other proposal 3 status = %s, want rejected", p3.Status())
	}
}

func TestGig_AcceptProposal_NotOpen(t *testing.T) {
	gig := createTestGig()
	gig.Cancel("test")

	_, err := gig.AcceptProposal(valueobject.GenerateProposalID(), valueobject.GenerateContractID())
	if err != ErrGigNotOpen {
		t.Errorf("AcceptProposal() on cancelled gig error = %v, want ErrGigNotOpen", err)
	}
}

func TestGig_FindProposal(t *testing.T) {
	gig := createTestGig()
	proposalID := valueobject.GenerateProposalID()
	price := valueobject.MustNewMoney(25000, valueobject.NGN)
	proposal := NewProposal(proposalID, valueobject.GenerateUserID(), "Test", price, 7, nil)
	gig.SubmitProposal(proposal)

	found := gig.FindProposal(proposalID)
	if found == nil {
		t.Fatal("FindProposal() should return proposal")
	}
	if !found.ID().Equals(proposalID) {
		t.Error("FindProposal() should return correct proposal")
	}

	notFound := gig.FindProposal(valueobject.GenerateProposalID())
	if notFound != nil {
		t.Error("FindProposal() should return nil for non-existent")
	}
}

func TestGig_FindProposalByHustler(t *testing.T) {
	gig := createTestGig()
	hustlerID := valueobject.GenerateUserID()
	price := valueobject.MustNewMoney(25000, valueobject.NGN)
	proposal := NewProposal(valueobject.GenerateProposalID(), hustlerID, "Test", price, 7, nil)
	gig.SubmitProposal(proposal)

	found := gig.FindProposalByHustler(hustlerID)
	if found == nil {
		t.Fatal("FindProposalByHustler() should return proposal")
	}
	if !found.HustlerID().Equals(hustlerID) {
		t.Error("FindProposalByHustler() should return correct proposal")
	}
}

func TestGig_MarkCompleted(t *testing.T) {
	gig := createTestGig()

	gig.MarkCompleted()

	if gig.Status() != GigStatusCompleted {
		t.Errorf("MarkCompleted() status = %s, want completed", gig.Status())
	}
}

func TestGig_MarkDisputed(t *testing.T) {
	gig := createTestGig()

	gig.MarkDisputed()

	if gig.Status() != GigStatusDisputed {
		t.Errorf("MarkDisputed() status = %s, want disputed", gig.Status())
	}
}

func TestGigStatus_String(t *testing.T) {
	tests := []struct {
		status GigStatus
		want   string
	}{
		{GigStatusOpen, "open"},
		{GigStatusInProgress, "in_progress"},
		{GigStatusCompleted, "completed"},
		{GigStatusCancelled, "cancelled"},
		{GigStatusDisputed, "disputed"},
	}

	for _, tt := range tests {
		if got := tt.status.String(); got != tt.want {
			t.Errorf("%v.String() = %s, want %s", tt.status, got, tt.want)
		}
	}
}

func TestReconstructGig(t *testing.T) {
	gigID := valueobject.GenerateGigID()
	clientID := valueobject.GenerateUserID()
	budget := createTestBudget()
	now := time.Now().UTC()

	gig := ReconstructGig(
		gigID,
		clientID,
		"Test Gig",
		"Description",
		"Category",
		nil,
		budget,
		valueobject.NGN,
		14,
		nil,
		true,
		"",
		GigStatusInProgress,
		10,
		true,
		[]string{"file.pdf"},
		[]string{"tag1"},
		[]*Proposal{},
		nil,
		now,
		now,
		5,
	)

	if !gig.ID().Equals(gigID) {
		t.Error("ReconstructGig() should set ID")
	}
	if gig.Status() != GigStatusInProgress {
		t.Errorf("ReconstructGig() status = %s", gig.Status())
	}
	if gig.ViewCount() != 10 {
		t.Errorf("ReconstructGig() viewCount = %d, want 10", gig.ViewCount())
	}
	if !gig.IsFeatured() {
		t.Error("ReconstructGig() should be featured")
	}
	if gig.Version() != 5 {
		t.Errorf("ReconstructGig() version = %d, want 5", gig.Version())
	}

	// Should not record any events
	if len(gig.DomainEvents()) != 0 {
		t.Error("ReconstructGig() should not record events")
	}
}

// Helper functions

func createTestBudget() Budget {
	min := valueobject.MustNewMoney(10000, valueobject.NGN)
	max := valueobject.MustNewMoney(50000, valueobject.NGN)
	budget, _ := NewBudget(min, max)
	return budget
}

func createTestProposal() *Proposal {
	return NewProposal(
		valueobject.GenerateProposalID(),
		valueobject.GenerateUserID(),
		"Test cover letter",
		valueobject.MustNewMoney(25000, valueobject.NGN),
		7,
		nil,
	)
}

func createTestProposalWithPrice(amount int64) *Proposal {
	return NewProposal(
		valueobject.GenerateProposalID(),
		valueobject.GenerateUserID(),
		"Test cover letter",
		valueobject.MustNewMoney(amount, valueobject.NGN),
		7,
		nil,
	)
}

func createTestGig() *Gig {
	gig, _ := NewGig(
		valueobject.GenerateGigID(),
		valueobject.GenerateUserID(),
		"Test Gig",
		"Test Description",
		"Development",
		createTestBudget(),
		14,
		true,
	)
	return gig
}

func createTestGigWithClient(clientID valueobject.UserID) *Gig {
	gig, _ := NewGig(
		valueobject.GenerateGigID(),
		clientID,
		"Test Gig",
		"Test Description",
		"Development",
		createTestBudget(),
		14,
		true,
	)
	return gig
}
