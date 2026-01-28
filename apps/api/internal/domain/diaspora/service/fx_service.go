package service

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"hustlex/internal/domain/diaspora/entity"
)

// Corridor represents a remittance corridor with rates and limits
type Corridor struct {
	SourceCurrency entity.SupportedCurrency
	TargetCurrency entity.SupportedCurrency
	SpreadBps      int             // Spread in basis points (100 bps = 1%)
	MinAmount      decimal.Decimal
	MaxAmount      decimal.Decimal
	FixedFee       decimal.Decimal
	PercentFee     decimal.Decimal // As decimal (0.015 = 1.5%)
	IsActive       bool
	DeliveryETA    string
}

// FXQuote represents a foreign exchange quote
type FXQuote struct {
	ID              string
	SourceCurrency  entity.SupportedCurrency
	TargetCurrency  entity.SupportedCurrency
	MidRate         decimal.Decimal
	BuyRate         decimal.Decimal // Rate for buying target currency
	SellRate        decimal.Decimal // Rate for selling target currency
	SpreadBps       int
	SourceAmount    decimal.Decimal
	TargetAmount    decimal.Decimal
	Fee             decimal.Decimal
	TotalSource     decimal.Decimal
	ValidUntil      time.Time
	CreatedAt       time.Time
}

// FXRateProvider interface for getting live rates
type FXRateProvider interface {
	GetRate(ctx context.Context, source, target entity.SupportedCurrency) (decimal.Decimal, error)
}

// FXService handles foreign exchange operations
type FXService struct {
	corridors    map[string]*Corridor // key: "SOURCE_TARGET"
	rateCache    map[string]*cachedRate
	cacheMu      sync.RWMutex
	rateProvider FXRateProvider
	quoteTTL     time.Duration
}

type cachedRate struct {
	rate      decimal.Decimal
	expiresAt time.Time
}

// NewFXService creates a new FX service
func NewFXService(provider FXRateProvider) *FXService {
	svc := &FXService{
		corridors:    make(map[string]*Corridor),
		rateCache:    make(map[string]*cachedRate),
		rateProvider: provider,
		quoteTTL:     15 * time.Minute, // Quotes valid for 15 minutes
	}

	// Initialize default corridors
	svc.initializeCorridors()

	return svc
}

// initializeCorridors sets up the supported remittance corridors
func (s *FXService) initializeCorridors() {
	// GBP to NGN corridor (UK to Nigeria)
	s.corridors["GBP_NGN"] = &Corridor{
		SourceCurrency: entity.CurrencyGBP,
		TargetCurrency: entity.CurrencyNGN,
		SpreadBps:      150, // 1.5% spread
		MinAmount:      decimal.NewFromInt(10),
		MaxAmount:      decimal.NewFromInt(10000),
		FixedFee:       decimal.NewFromFloat(2.99),
		PercentFee:     decimal.NewFromFloat(0.005), // 0.5%
		IsActive:       true,
		DeliveryETA:    "Within 24 hours",
	}

	// USD to NGN corridor (US to Nigeria)
	s.corridors["USD_NGN"] = &Corridor{
		SourceCurrency: entity.CurrencyUSD,
		TargetCurrency: entity.CurrencyNGN,
		SpreadBps:      175, // 1.75% spread
		MinAmount:      decimal.NewFromInt(10),
		MaxAmount:      decimal.NewFromInt(10000),
		FixedFee:       decimal.NewFromFloat(2.99),
		PercentFee:     decimal.NewFromFloat(0.005),
		IsActive:       true,
		DeliveryETA:    "Within 24 hours",
	}

	// EUR to NGN corridor (Europe to Nigeria)
	s.corridors["EUR_NGN"] = &Corridor{
		SourceCurrency: entity.CurrencyEUR,
		TargetCurrency: entity.CurrencyNGN,
		SpreadBps:      175,
		MinAmount:      decimal.NewFromInt(10),
		MaxAmount:      decimal.NewFromInt(10000),
		FixedFee:       decimal.NewFromFloat(2.99),
		PercentFee:     decimal.NewFromFloat(0.005),
		IsActive:       true,
		DeliveryETA:    "Within 24 hours",
	}

	// CAD to NGN corridor (Canada to Nigeria)
	s.corridors["CAD_NGN"] = &Corridor{
		SourceCurrency: entity.CurrencyCAD,
		TargetCurrency: entity.CurrencyNGN,
		SpreadBps:      200, // 2% spread
		MinAmount:      decimal.NewFromInt(10),
		MaxAmount:      decimal.NewFromInt(10000),
		FixedFee:       decimal.NewFromFloat(3.99),
		PercentFee:     decimal.NewFromFloat(0.005),
		IsActive:       true,
		DeliveryETA:    "Within 24 hours",
	}

	// GBP to GHS corridor (UK to Ghana)
	s.corridors["GBP_GHS"] = &Corridor{
		SourceCurrency: entity.CurrencyGBP,
		TargetCurrency: entity.CurrencyGHS,
		SpreadBps:      200,
		MinAmount:      decimal.NewFromInt(10),
		MaxAmount:      decimal.NewFromInt(5000),
		FixedFee:       decimal.NewFromFloat(3.99),
		PercentFee:     decimal.NewFromFloat(0.01), // 1%
		IsActive:       true,
		DeliveryETA:    "Within 24 hours",
	}

	// USD to GHS corridor
	s.corridors["USD_GHS"] = &Corridor{
		SourceCurrency: entity.CurrencyUSD,
		TargetCurrency: entity.CurrencyGHS,
		SpreadBps:      200,
		MinAmount:      decimal.NewFromInt(10),
		MaxAmount:      decimal.NewFromInt(5000),
		FixedFee:       decimal.NewFromFloat(3.99),
		PercentFee:     decimal.NewFromFloat(0.01),
		IsActive:       true,
		DeliveryETA:    "Within 24 hours",
	}

	// GBP to KES corridor (UK to Kenya)
	s.corridors["GBP_KES"] = &Corridor{
		SourceCurrency: entity.CurrencyGBP,
		TargetCurrency: entity.CurrencyKES,
		SpreadBps:      225, // 2.25% spread
		MinAmount:      decimal.NewFromInt(10),
		MaxAmount:      decimal.NewFromInt(5000),
		FixedFee:       decimal.NewFromFloat(3.99),
		PercentFee:     decimal.NewFromFloat(0.01),
		IsActive:       true,
		DeliveryETA:    "1-2 business days",
	}

	// USD to KES corridor
	s.corridors["USD_KES"] = &Corridor{
		SourceCurrency: entity.CurrencyUSD,
		TargetCurrency: entity.CurrencyKES,
		SpreadBps:      225,
		MinAmount:      decimal.NewFromInt(10),
		MaxAmount:      decimal.NewFromInt(5000),
		FixedFee:       decimal.NewFromFloat(3.99),
		PercentFee:     decimal.NewFromFloat(0.01),
		IsActive:       true,
		DeliveryETA:    "1-2 business days",
	}

	// Internal NGN transfers (no FX)
	s.corridors["NGN_NGN"] = &Corridor{
		SourceCurrency: entity.CurrencyNGN,
		TargetCurrency: entity.CurrencyNGN,
		SpreadBps:      0,
		MinAmount:      decimal.NewFromInt(100),
		MaxAmount:      decimal.NewFromInt(5000000),
		FixedFee:       decimal.NewFromInt(50),
		PercentFee:     decimal.Zero,
		IsActive:       true,
		DeliveryETA:    "Instant",
	}
}

// GetCorridor returns a corridor if it exists and is active
func (s *FXService) GetCorridor(source, target entity.SupportedCurrency) (*Corridor, error) {
	key := string(source) + "_" + string(target)
	corridor, exists := s.corridors[key]
	if !exists {
		return nil, errors.New("corridor not supported")
	}
	if !corridor.IsActive {
		return nil, errors.New("corridor temporarily unavailable")
	}
	return corridor, nil
}

// GetAllCorridors returns all active corridors
func (s *FXService) GetAllCorridors() []*Corridor {
	corridors := make([]*Corridor, 0, len(s.corridors))
	for _, c := range s.corridors {
		if c.IsActive {
			corridors = append(corridors, c)
		}
	}
	return corridors
}

// GetQuote generates an FX quote for a transfer
func (s *FXService) GetQuote(ctx context.Context, source, target entity.SupportedCurrency, sourceAmount decimal.Decimal) (*FXQuote, error) {
	corridor, err := s.GetCorridor(source, target)
	if err != nil {
		return nil, err
	}

	// Validate amount
	if sourceAmount.LessThan(corridor.MinAmount) {
		return nil, errors.New("amount below minimum")
	}
	if sourceAmount.GreaterThan(corridor.MaxAmount) {
		return nil, errors.New("amount exceeds maximum")
	}

	// Get mid-market rate
	midRate, err := s.getMidRate(ctx, source, target)
	if err != nil {
		return nil, err
	}

	// Calculate spread
	spreadMultiplier := decimal.NewFromInt(int64(corridor.SpreadBps)).Div(decimal.NewFromInt(10000))

	// Buy rate (customer buys target currency) - less favorable
	buyRate := midRate.Mul(decimal.NewFromInt(1).Sub(spreadMultiplier))

	// Sell rate (customer sells target currency) - more favorable
	sellRate := midRate.Mul(decimal.NewFromInt(1).Add(spreadMultiplier))

	// Calculate fee
	percentFee := sourceAmount.Mul(corridor.PercentFee)
	totalFee := corridor.FixedFee.Add(percentFee)

	// Amount after fee
	netSourceAmount := sourceAmount.Sub(totalFee)

	// Calculate target amount using buy rate
	targetAmount := netSourceAmount.Mul(buyRate)

	now := time.Now()
	quote := &FXQuote{
		ID:             uuid.New().String(),
		SourceCurrency: source,
		TargetCurrency: target,
		MidRate:        midRate,
		BuyRate:        buyRate,
		SellRate:       sellRate,
		SpreadBps:      corridor.SpreadBps,
		SourceAmount:   sourceAmount,
		TargetAmount:   targetAmount.Round(2),
		Fee:            totalFee.Round(2),
		TotalSource:    sourceAmount,
		ValidUntil:     now.Add(s.quoteTTL),
		CreatedAt:      now,
	}

	return quote, nil
}

// GetQuoteByTargetAmount generates a quote when target amount is specified
func (s *FXService) GetQuoteByTargetAmount(ctx context.Context, source, target entity.SupportedCurrency, targetAmount decimal.Decimal) (*FXQuote, error) {
	corridor, err := s.GetCorridor(source, target)
	if err != nil {
		return nil, err
	}

	// Get mid-market rate
	midRate, err := s.getMidRate(ctx, source, target)
	if err != nil {
		return nil, err
	}

	// Calculate spread
	spreadMultiplier := decimal.NewFromInt(int64(corridor.SpreadBps)).Div(decimal.NewFromInt(10000))
	buyRate := midRate.Mul(decimal.NewFromInt(1).Sub(spreadMultiplier))
	sellRate := midRate.Mul(decimal.NewFromInt(1).Add(spreadMultiplier))

	// Calculate source amount needed
	netSourceAmount := targetAmount.Div(buyRate)

	// Add fees to get total source
	// netSourceAmount = sourceAmount - (fixedFee + sourceAmount * percentFee)
	// netSourceAmount = sourceAmount * (1 - percentFee) - fixedFee
	// sourceAmount = (netSourceAmount + fixedFee) / (1 - percentFee)
	denominator := decimal.NewFromInt(1).Sub(corridor.PercentFee)
	sourceAmount := netSourceAmount.Add(corridor.FixedFee).Div(denominator)

	// Validate amount
	if sourceAmount.LessThan(corridor.MinAmount) {
		return nil, errors.New("calculated amount below minimum")
	}
	if sourceAmount.GreaterThan(corridor.MaxAmount) {
		return nil, errors.New("calculated amount exceeds maximum")
	}

	percentFee := sourceAmount.Mul(corridor.PercentFee)
	totalFee := corridor.FixedFee.Add(percentFee)

	now := time.Now()
	quote := &FXQuote{
		ID:             uuid.New().String(),
		SourceCurrency: source,
		TargetCurrency: target,
		MidRate:        midRate,
		BuyRate:        buyRate,
		SellRate:       sellRate,
		SpreadBps:      corridor.SpreadBps,
		SourceAmount:   sourceAmount.Round(2),
		TargetAmount:   targetAmount,
		Fee:            totalFee.Round(2),
		TotalSource:    sourceAmount.Round(2),
		ValidUntil:     now.Add(s.quoteTTL),
		CreatedAt:      now,
	}

	return quote, nil
}

// getMidRate gets the mid-market rate (from cache or provider)
func (s *FXService) getMidRate(ctx context.Context, source, target entity.SupportedCurrency) (decimal.Decimal, error) {
	// Same currency
	if source == target {
		return decimal.NewFromInt(1), nil
	}

	key := string(source) + "_" + string(target)

	// Check cache
	s.cacheMu.RLock()
	cached, exists := s.rateCache[key]
	s.cacheMu.RUnlock()

	if exists && time.Now().Before(cached.expiresAt) {
		return cached.rate, nil
	}

	// Get from provider
	if s.rateProvider != nil {
		rate, err := s.rateProvider.GetRate(ctx, source, target)
		if err == nil {
			// Cache the rate
			s.cacheMu.Lock()
			s.rateCache[key] = &cachedRate{
				rate:      rate,
				expiresAt: time.Now().Add(5 * time.Minute),
			}
			s.cacheMu.Unlock()
			return rate, nil
		}
	}

	// Fallback to static rates (for development/testing)
	return s.getStaticRate(source, target)
}

// getStaticRate returns static fallback rates
func (s *FXService) getStaticRate(source, target entity.SupportedCurrency) (decimal.Decimal, error) {
	// Static rates to NGN (approximate market rates)
	toNGN := map[entity.SupportedCurrency]decimal.Decimal{
		entity.CurrencyGBP: decimal.NewFromFloat(1950.00),
		entity.CurrencyUSD: decimal.NewFromFloat(1550.00),
		entity.CurrencyEUR: decimal.NewFromFloat(1700.00),
		entity.CurrencyCAD: decimal.NewFromFloat(1150.00),
		entity.CurrencyGHS: decimal.NewFromFloat(130.00),
		entity.CurrencyKES: decimal.NewFromFloat(12.00),
		entity.CurrencyNGN: decimal.NewFromInt(1),
	}

	if target == entity.CurrencyNGN {
		if rate, ok := toNGN[source]; ok {
			return rate, nil
		}
	}

	// Calculate cross rates via NGN
	sourceToNGN, sourceOK := toNGN[source]
	targetToNGN, targetOK := toNGN[target]

	if sourceOK && targetOK {
		return sourceToNGN.Div(targetToNGN), nil
	}

	return decimal.Zero, errors.New("rate not available")
}

// ValidateQuote checks if a quote is still valid
func (s *FXService) ValidateQuote(quote *FXQuote) bool {
	return time.Now().Before(quote.ValidUntil)
}

// CalculateFee calculates the fee for a transfer
func (s *FXService) CalculateFee(source, target entity.SupportedCurrency, amount decimal.Decimal) (decimal.Decimal, error) {
	corridor, err := s.GetCorridor(source, target)
	if err != nil {
		return decimal.Zero, err
	}

	percentFee := amount.Mul(corridor.PercentFee)
	return corridor.FixedFee.Add(percentFee).Round(2), nil
}
