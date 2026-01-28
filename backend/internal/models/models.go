package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserTier represents the credit tier of a user
type UserTier string

const (
	TierBronze   UserTier = "bronze"
	TierSilver   UserTier = "silver"
	TierGold     UserTier = "gold"
	TierPlatinum UserTier = "platinum"
)

// GigStatus represents the status of a gig
type GigStatus string

const (
	GigStatusOpen       GigStatus = "open"
	GigStatusInProgress GigStatus = "in_progress"
	GigStatusCompleted  GigStatus = "completed"
	GigStatusCancelled  GigStatus = "cancelled"
	GigStatusDisputed   GigStatus = "disputed"
)

// ContractStatus represents the status of a gig contract
type ContractStatus string

const (
	ContractStatusActive    ContractStatus = "active"
	ContractStatusDelivered ContractStatus = "delivered"
	ContractStatusCompleted ContractStatus = "completed"
	ContractStatusDisputed  ContractStatus = "disputed"
	ContractStatusCancelled ContractStatus = "cancelled"
)

// CircleType represents the type of savings circle
type CircleType string

const (
	CircleTypeRotational  CircleType = "rotational"
	CircleTypeFixedTarget CircleType = "fixed_target"
	CircleTypeEmergency   CircleType = "emergency"
)

// ContributionStatus represents the status of a contribution
type ContributionStatus string

const (
	ContributionStatusPending ContributionStatus = "pending"
	ContributionStatusPaid    ContributionStatus = "paid"
	ContributionStatusOverdue ContributionStatus = "overdue"
	ContributionStatusWaived  ContributionStatus = "waived"
)

// TransactionType represents the type of wallet transaction
type TransactionType string

const (
	TransactionTypeDeposit        TransactionType = "deposit"
	TransactionTypeWithdrawal     TransactionType = "withdrawal"
	TransactionTypeGigPayment     TransactionType = "gig_payment"
	TransactionTypeGigEarning     TransactionType = "gig_earning"
	TransactionTypeContribution   TransactionType = "contribution"
	TransactionTypePayout         TransactionType = "payout"
	TransactionTypeTransferIn     TransactionType = "transfer_in"
	TransactionTypeTransferOut    TransactionType = "transfer_out"
	TransactionTypeLoanDisburse   TransactionType = "loan_disburse"
	TransactionTypeLoanRepayment  TransactionType = "loan_repayment"
	TransactionTypeRefund         TransactionType = "refund"
	TransactionTypeFee            TransactionType = "fee"
)

// TransactionStatus represents the status of a transaction
type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "pending"
	TransactionStatusCompleted TransactionStatus = "completed"
	TransactionStatusFailed    TransactionStatus = "failed"
	TransactionStatusReversed  TransactionStatus = "reversed"
)

// BaseModel contains common fields for all models
type BaseModel struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// User represents a HustleX user
type User struct {
	BaseModel
	Phone           string    `gorm:"uniqueIndex;not null;size:15" json:"phone"`
	Email           string    `gorm:"uniqueIndex;size:255" json:"email,omitempty"`
	FullName        string    `gorm:"size:255;not null" json:"full_name"`
	Username        string    `gorm:"uniqueIndex;size:50" json:"username,omitempty"`
	ProfileImage    string    `gorm:"size:500" json:"profile_image,omitempty"`
	Bio             string    `gorm:"size:500" json:"bio,omitempty"`
	Location        string    `gorm:"size:100" json:"location,omitempty"`
	State           string    `gorm:"size:50" json:"state,omitempty"`
	DateOfBirth     *time.Time `json:"date_of_birth,omitempty"`
	Gender          string    `gorm:"size:20" json:"gender,omitempty"`
	IsVerified      bool      `gorm:"default:false" json:"is_verified"`
	IsActive        bool      `gorm:"default:true" json:"is_active"`
	Tier            UserTier  `gorm:"type:varchar(20);default:'bronze'" json:"tier"`
	ReferralCode    string    `gorm:"uniqueIndex;size:10" json:"referral_code"`
	ReferredBy      *uuid.UUID `gorm:"type:uuid" json:"referred_by,omitempty"`
	LastLoginAt     *time.Time `json:"last_login_at,omitempty"`
	
	// Relationships
	Skills        []UserSkill     `gorm:"foreignKey:UserID" json:"skills,omitempty"`
	Wallet        *Wallet         `gorm:"foreignKey:UserID" json:"wallet,omitempty"`
	CreditScore   *CreditScore    `gorm:"foreignKey:UserID" json:"credit_score,omitempty"`
	GigsPosted    []Gig           `gorm:"foreignKey:ClientID" json:"-"`
	Proposals     []GigProposal   `gorm:"foreignKey:HustlerID" json:"-"`
	CircleMembers []CircleMember  `gorm:"foreignKey:UserID" json:"-"`
}

// Skill represents a skill category
type Skill struct {
	BaseModel
	Name        string `gorm:"uniqueIndex;size:100;not null" json:"name"`
	Category    string `gorm:"size:50;not null;index" json:"category"`
	Description string `gorm:"size:500" json:"description,omitempty"`
	Icon        string `gorm:"size:100" json:"icon,omitempty"`
	IsActive    bool   `gorm:"default:true" json:"is_active"`
}

// UserSkill represents a user's skill with proficiency level
type UserSkill struct {
	BaseModel
	UserID        uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	SkillID       uuid.UUID `gorm:"type:uuid;not null;index" json:"skill_id"`
	Proficiency   string    `gorm:"size:20;default:'intermediate'" json:"proficiency"` // beginner, intermediate, expert
	YearsExp      int       `gorm:"default:0" json:"years_experience"`
	IsVerified    bool      `gorm:"default:false" json:"is_verified"`
	PortfolioURLs []string  `gorm:"type:text[];serializer:json" json:"portfolio_urls,omitempty"`
	
	// Relationships
	User  User  `gorm:"foreignKey:UserID" json:"-"`
	Skill Skill `gorm:"foreignKey:SkillID" json:"skill,omitempty"`
}

// Gig represents a job/task posted by a client
type Gig struct {
	BaseModel
	ClientID     uuid.UUID  `gorm:"type:uuid;not null;index" json:"client_id"`
	Title        string     `gorm:"size:200;not null" json:"title"`
	Description  string     `gorm:"type:text;not null" json:"description"`
	Category     string     `gorm:"size:50;not null;index" json:"category"`
	SkillID      *uuid.UUID `gorm:"type:uuid;index" json:"skill_id,omitempty"`
	BudgetMin    int64      `gorm:"not null" json:"budget_min"` // in kobo (smallest unit)
	BudgetMax    int64      `gorm:"not null" json:"budget_max"` // in kobo
	Currency     string     `gorm:"size:3;default:'NGN'" json:"currency"`
	Deadline     *time.Time `json:"deadline,omitempty"`
	DeliveryDays int        `gorm:"default:7" json:"delivery_days"`
	IsRemote     bool       `gorm:"default:true" json:"is_remote"`
	Location     string     `gorm:"size:100" json:"location,omitempty"`
	Status       GigStatus  `gorm:"type:varchar(20);default:'open';index" json:"status"`
	ViewCount    int        `gorm:"default:0" json:"view_count"`
	ProposalCount int       `gorm:"default:0" json:"proposal_count"`
	IsFeatured   bool       `gorm:"default:false;index" json:"is_featured"`
	Attachments  []string   `gorm:"type:text[];serializer:json" json:"attachments,omitempty"`
	Tags         []string   `gorm:"type:text[];serializer:json" json:"tags,omitempty"`
	
	// Relationships
	Client    User          `gorm:"foreignKey:ClientID" json:"client,omitempty"`
	Skill     *Skill        `gorm:"foreignKey:SkillID" json:"skill,omitempty"`
	Proposals []GigProposal `gorm:"foreignKey:GigID" json:"proposals,omitempty"`
	Contract  *GigContract  `gorm:"foreignKey:GigID" json:"contract,omitempty"`
}

// GigProposal represents a proposal submitted by a hustler for a gig
type GigProposal struct {
	BaseModel
	GigID         uuid.UUID `gorm:"type:uuid;not null;index" json:"gig_id"`
	HustlerID     uuid.UUID `gorm:"type:uuid;not null;index" json:"hustler_id"`
	CoverLetter   string    `gorm:"type:text;not null" json:"cover_letter"`
	ProposedPrice int64     `gorm:"not null" json:"proposed_price"` // in kobo
	DeliveryDays  int       `gorm:"not null" json:"delivery_days"`
	Status        string    `gorm:"size:20;default:'pending';index" json:"status"` // pending, accepted, rejected, withdrawn
	Attachments   []string  `gorm:"type:text[];serializer:json" json:"attachments,omitempty"`
	
	// Relationships
	Gig     Gig  `gorm:"foreignKey:GigID" json:"gig,omitempty"`
	Hustler User `gorm:"foreignKey:HustlerID" json:"hustler,omitempty"`
}

// GigContract represents an accepted gig contract
type GigContract struct {
	BaseModel
	GigID        uuid.UUID      `gorm:"type:uuid;not null;uniqueIndex" json:"gig_id"`
	HustlerID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"hustler_id"`
	ProposalID   uuid.UUID      `gorm:"type:uuid;not null" json:"proposal_id"`
	AgreedPrice  int64          `gorm:"not null" json:"agreed_price"` // in kobo
	PlatformFee  int64          `gorm:"not null" json:"platform_fee"` // 10% fee
	DeliveryDays int            `gorm:"not null" json:"delivery_days"`
	Status       ContractStatus `gorm:"type:varchar(20);default:'active';index" json:"status"`
	StartedAt    time.Time      `gorm:"autoCreateTime" json:"started_at"`
	DeadlineAt   time.Time      `json:"deadline_at"`
	DeliveredAt  *time.Time     `json:"delivered_at,omitempty"`
	CompletedAt  *time.Time     `json:"completed_at,omitempty"`
	Deliverables []string       `gorm:"type:text[];serializer:json" json:"deliverables,omitempty"`
	ClientNotes  string         `gorm:"type:text" json:"client_notes,omitempty"`
	
	// Relationships
	Gig      Gig          `gorm:"foreignKey:GigID" json:"gig,omitempty"`
	Hustler  User         `gorm:"foreignKey:HustlerID" json:"hustler,omitempty"`
	Proposal GigProposal  `gorm:"foreignKey:ProposalID" json:"proposal,omitempty"`
	Review   *GigReview   `gorm:"foreignKey:ContractID" json:"review,omitempty"`
}

// GigReview represents a review for a completed gig
type GigReview struct {
	BaseModel
	ContractID  uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"contract_id"`
	ReviewerID  uuid.UUID `gorm:"type:uuid;not null;index" json:"reviewer_id"`
	RevieweeID  uuid.UUID `gorm:"type:uuid;not null;index" json:"reviewee_id"`
	Rating      int       `gorm:"not null;check:rating >= 1 AND rating <= 5" json:"rating"` // 1-5 stars
	ReviewText  string    `gorm:"type:text" json:"review_text,omitempty"`
	IsPublic    bool      `gorm:"default:true" json:"is_public"`
	
	// Detailed ratings
	CommunicationRating int `gorm:"check:communication_rating >= 1 AND communication_rating <= 5" json:"communication_rating,omitempty"`
	QualityRating       int `gorm:"check:quality_rating >= 1 AND quality_rating <= 5" json:"quality_rating,omitempty"`
	TimelinessRating    int `gorm:"check:timeliness_rating >= 1 AND timeliness_rating <= 5" json:"timeliness_rating,omitempty"`
	
	// Relationships
	Contract GigContract `gorm:"foreignKey:ContractID" json:"-"`
	Reviewer User        `gorm:"foreignKey:ReviewerID" json:"reviewer,omitempty"`
	Reviewee User        `gorm:"foreignKey:RevieweeID" json:"reviewee,omitempty"`
}

// SavingsCircle represents a group savings (Ajo/Esusu) circle
type SavingsCircle struct {
	BaseModel
	Name            string     `gorm:"size:100;not null" json:"name"`
	Description     string     `gorm:"size:500" json:"description,omitempty"`
	Type            CircleType `gorm:"type:varchar(20);not null" json:"type"`
	ContributionAmt int64      `gorm:"not null" json:"contribution_amount"` // in kobo
	Currency        string     `gorm:"size:3;default:'NGN'" json:"currency"`
	Frequency       string     `gorm:"size:20;not null" json:"frequency"` // daily, weekly, biweekly, monthly
	MaxMembers      int        `gorm:"not null;default:12" json:"max_members"`
	CurrentMembers  int        `gorm:"default:1" json:"current_members"`
	CurrentRound    int        `gorm:"default:0" json:"current_round"`
	TotalRounds     int        `gorm:"not null" json:"total_rounds"`
	PoolBalance     int64      `gorm:"default:0" json:"pool_balance"` // current round's pool
	TotalSaved      int64      `gorm:"default:0" json:"total_saved"` // lifetime saved
	CreatedBy       uuid.UUID  `gorm:"type:uuid;not null;index" json:"created_by"`
	Status          string     `gorm:"size:20;default:'recruiting';index" json:"status"` // recruiting, active, completed, cancelled
	StartDate       *time.Time `json:"start_date,omitempty"`
	NextPayoutDate  *time.Time `json:"next_payout_date,omitempty"`
	IsPrivate       bool       `gorm:"default:false" json:"is_private"`
	InviteCode      string     `gorm:"uniqueIndex;size:10" json:"invite_code"`
	Rules           []string   `gorm:"type:text[];serializer:json" json:"rules,omitempty"`
	
	// Relationships
	Creator       User           `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	Members       []CircleMember `gorm:"foreignKey:CircleID" json:"members,omitempty"`
	Contributions []Contribution `gorm:"foreignKey:CircleID" json:"-"`
}

// CircleMember represents a member of a savings circle
type CircleMember struct {
	BaseModel
	CircleID      uuid.UUID `gorm:"type:uuid;not null;index" json:"circle_id"`
	UserID        uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Position      int       `gorm:"not null" json:"position"` // payout position in rotation
	Role          string    `gorm:"size:20;default:'member'" json:"role"` // admin, member
	Status        string    `gorm:"size:20;default:'active'" json:"status"` // pending, active, removed, left
	JoinedAt      time.Time `gorm:"autoCreateTime" json:"joined_at"`
	TotalContrib  int64     `gorm:"default:0" json:"total_contributed"` // lifetime contributions
	MissedPayments int      `gorm:"default:0" json:"missed_payments"`
	HasReceived   bool      `gorm:"default:false" json:"has_received"` // received payout this cycle
	
	// Relationships
	Circle Circle `gorm:"foreignKey:CircleID" json:"-"`
	User   User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// Contribution represents a single contribution to a savings circle
type Contribution struct {
	BaseModel
	CircleID      uuid.UUID          `gorm:"type:uuid;not null;index" json:"circle_id"`
	MemberID      uuid.UUID          `gorm:"type:uuid;not null;index" json:"member_id"`
	Round         int                `gorm:"not null" json:"round"`
	Amount        int64              `gorm:"not null" json:"amount"` // in kobo
	DueDate       time.Time          `gorm:"not null" json:"due_date"`
	PaidAt        *time.Time         `json:"paid_at,omitempty"`
	Status        ContributionStatus `gorm:"type:varchar(20);default:'pending';index" json:"status"`
	TransactionID *uuid.UUID         `gorm:"type:uuid" json:"transaction_id,omitempty"`
	LateFee       int64              `gorm:"default:0" json:"late_fee"`
	
	// Relationships
	Circle      SavingsCircle `gorm:"foreignKey:CircleID" json:"-"`
	Member      CircleMember  `gorm:"foreignKey:MemberID" json:"member,omitempty"`
	Transaction *Transaction  `gorm:"foreignKey:TransactionID" json:"-"`
}

// Wallet represents a user's wallet
type Wallet struct {
	BaseModel
	UserID        uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`
	Balance       int64     `gorm:"default:0;not null" json:"balance"` // available balance in kobo
	EscrowBalance int64     `gorm:"default:0;not null" json:"escrow_balance"` // held for gigs
	SavingsBalance int64    `gorm:"default:0;not null" json:"savings_balance"` // in savings circles
	Currency      string    `gorm:"size:3;default:'NGN'" json:"currency"`
	IsLocked      bool      `gorm:"default:false" json:"is_locked"`
	Pin           string    `gorm:"size:255" json:"-"` // hashed transaction PIN
	PinAttempts   int       `gorm:"default:0" json:"-"`
	
	// Relationships
	User         User          `gorm:"foreignKey:UserID" json:"-"`
	Transactions []Transaction `gorm:"foreignKey:WalletID" json:"-"`
}

// Transaction represents a wallet transaction
type Transaction struct {
	BaseModel
	WalletID    uuid.UUID         `gorm:"type:uuid;not null;index" json:"wallet_id"`
	Type        TransactionType   `gorm:"type:varchar(30);not null;index" json:"type"`
	Amount      int64             `gorm:"not null" json:"amount"` // in kobo (positive)
	Fee         int64             `gorm:"default:0" json:"fee"` // platform fee
	NetAmount   int64             `gorm:"not null" json:"net_amount"` // amount after fee
	Currency    string            `gorm:"size:3;default:'NGN'" json:"currency"`
	Reference   string            `gorm:"uniqueIndex;size:50;not null" json:"reference"`
	ExternalRef string            `gorm:"size:100" json:"external_ref,omitempty"` // payment gateway ref
	Status      TransactionStatus `gorm:"type:varchar(20);default:'pending';index" json:"status"`
	Description string            `gorm:"size:255" json:"description,omitempty"`
	Metadata    map[string]interface{} `gorm:"type:jsonb;serializer:json" json:"metadata,omitempty"`
	BalanceBefore int64           `gorm:"not null" json:"balance_before"`
	BalanceAfter  int64           `gorm:"not null" json:"balance_after"`
	
	// For P2P transfers
	CounterpartyID *uuid.UUID `gorm:"type:uuid" json:"counterparty_id,omitempty"`
	
	// Relationships
	Wallet       Wallet `gorm:"foreignKey:WalletID" json:"-"`
	Counterparty *User  `gorm:"foreignKey:CounterpartyID" json:"counterparty,omitempty"`
}

// CreditScore represents a user's hustle credit score
type CreditScore struct {
	BaseModel
	UserID         uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`
	Score          int       `gorm:"not null;default:0;check:score >= 0 AND score <= 850" json:"score"`
	Tier           UserTier  `gorm:"type:varchar(20);default:'bronze'" json:"tier"`
	
	// Score components (0-100 each)
	GigCompletionScore   int `gorm:"default:0" json:"gig_completion_score"`
	RatingScore          int `gorm:"default:0" json:"rating_score"`
	SavingsScore         int `gorm:"default:0" json:"savings_score"`
	AccountAgeScore      int `gorm:"default:0" json:"account_age_score"`
	VerificationScore    int `gorm:"default:0" json:"verification_score"`
	CommunityScore       int `gorm:"default:0" json:"community_score"`
	
	// Stats
	TotalGigsCompleted   int     `gorm:"default:0" json:"total_gigs_completed"`
	TotalGigsAccepted    int     `gorm:"default:0" json:"total_gigs_accepted"`
	AverageRating        float64 `gorm:"default:0" json:"average_rating"`
	TotalReviews         int     `gorm:"default:0" json:"total_reviews"`
	OnTimeContributions  int     `gorm:"default:0" json:"on_time_contributions"`
	TotalContributions   int     `gorm:"default:0" json:"total_contributions"`
	
	LastCalculatedAt     time.Time `json:"last_calculated_at"`
	
	// Relationships
	User User `gorm:"foreignKey:UserID" json:"-"`
}

// Loan represents a micro-loan
type Loan struct {
	BaseModel
	UserID          uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Amount          int64     `gorm:"not null" json:"amount"` // principal in kobo
	InterestRate    float64   `gorm:"not null" json:"interest_rate"` // monthly rate (e.g., 0.05 for 5%)
	InterestAmount  int64     `gorm:"not null" json:"interest_amount"`
	TotalAmount     int64     `gorm:"not null" json:"total_amount"` // principal + interest
	AmountRepaid    int64     `gorm:"default:0" json:"amount_repaid"`
	Currency        string    `gorm:"size:3;default:'NGN'" json:"currency"`
	TenureMonths    int       `gorm:"not null" json:"tenure_months"`
	Status          string    `gorm:"size:20;default:'pending';index" json:"status"` // pending, approved, disbursed, repaying, completed, defaulted
	Purpose         string    `gorm:"size:255" json:"purpose,omitempty"`
	ApprovedAt      *time.Time `json:"approved_at,omitempty"`
	DisbursedAt     *time.Time `json:"disbursed_at,omitempty"`
	DueDate         *time.Time `json:"due_date,omitempty"`
	CompletedAt     *time.Time `json:"completed_at,omitempty"`
	
	// Relationships
	User       User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Repayments []LoanRepayment `gorm:"foreignKey:LoanID" json:"repayments,omitempty"`
}

// LoanRepayment represents a loan repayment
type LoanRepayment struct {
	BaseModel
	LoanID        uuid.UUID `gorm:"type:uuid;not null;index" json:"loan_id"`
	Amount        int64     `gorm:"not null" json:"amount"`
	TransactionID uuid.UUID `gorm:"type:uuid;not null" json:"transaction_id"`
	
	// Relationships
	Loan        Loan        `gorm:"foreignKey:LoanID" json:"-"`
	Transaction Transaction `gorm:"foreignKey:TransactionID" json:"-"`
}

// Course represents a learning course
type Course struct {
	BaseModel
	Title         string     `gorm:"size:200;not null" json:"title"`
	Description   string     `gorm:"type:text;not null" json:"description"`
	SkillID       *uuid.UUID `gorm:"type:uuid;index" json:"skill_id,omitempty"`
	Category      string     `gorm:"size:50;not null;index" json:"category"`
	Thumbnail     string     `gorm:"size:500" json:"thumbnail,omitempty"`
	DurationMins  int        `gorm:"not null" json:"duration_mins"`
	Difficulty    string     `gorm:"size:20;not null" json:"difficulty"` // beginner, intermediate, advanced
	InstructorID  *uuid.UUID `gorm:"type:uuid" json:"instructor_id,omitempty"`
	IsFree        bool       `gorm:"default:true" json:"is_free"`
	Price         int64      `gorm:"default:0" json:"price"` // in kobo if not free
	EnrollCount   int        `gorm:"default:0" json:"enroll_count"`
	Rating        float64    `gorm:"default:0" json:"rating"`
	RatingCount   int        `gorm:"default:0" json:"rating_count"`
	IsPublished   bool       `gorm:"default:false;index" json:"is_published"`
	Modules       []Module   `gorm:"type:jsonb;serializer:json" json:"modules,omitempty"`
	
	// Relationships
	Skill      *Skill `gorm:"foreignKey:SkillID" json:"skill,omitempty"`
	Instructor *User  `gorm:"foreignKey:InstructorID" json:"instructor,omitempty"`
}

// Module represents a course module
type Module struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	VideoURL    string `json:"video_url,omitempty"`
	ContentURL  string `json:"content_url,omitempty"`
	DurationMin int    `json:"duration_mins"`
	Order       int    `json:"order"`
}

// Enrollment represents a user's course enrollment
type Enrollment struct {
	BaseModel
	UserID       uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	CourseID     uuid.UUID  `gorm:"type:uuid;not null;index" json:"course_id"`
	Progress     int        `gorm:"default:0" json:"progress"` // percentage 0-100
	CompletedAt  *time.Time `json:"completed_at,omitempty"`
	Certificate  string     `gorm:"size:500" json:"certificate,omitempty"` // certificate URL
	LastAccessAt *time.Time `json:"last_access_at,omitempty"`
	
	// Relationships
	User   User   `gorm:"foreignKey:UserID" json:"-"`
	Course Course `gorm:"foreignKey:CourseID" json:"course,omitempty"`
}

// OTPCode represents a one-time password for authentication
type OTPCode struct {
	BaseModel
	Phone     string    `gorm:"index;not null;size:15" json:"phone"`
	Code      string    `gorm:"not null;size:10" json:"-"`
	Purpose   string    `gorm:"not null;size:20" json:"purpose"` // login, register, reset_pin
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	IsUsed    bool      `gorm:"default:false" json:"is_used"`
	Attempts  int       `gorm:"default:0" json:"-"`
}

// Notification represents a user notification
type Notification struct {
	BaseModel
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Title     string    `gorm:"size:200;not null" json:"title"`
	Body      string    `gorm:"type:text;not null" json:"body"`
	Type      string    `gorm:"size:50;not null;index" json:"type"` // gig, savings, wallet, system
	Data      map[string]interface{} `gorm:"type:jsonb;serializer:json" json:"data,omitempty"`
	IsRead    bool      `gorm:"default:false;index" json:"is_read"`
	ReadAt    *time.Time `json:"read_at,omitempty"`
	
	// Relationships
	User User `gorm:"foreignKey:UserID" json:"-"`
}

// Alias for backwards compatibility
type Circle = SavingsCircle
