package valueobject

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidPercentage = errors.New("percentage must be between 0 and 100")
)

// Percentage represents a percentage value (0-100)
type Percentage struct {
	value float64
}

// NewPercentage creates a new Percentage value object
func NewPercentage(value float64) (Percentage, error) {
	if value < 0 || value > 100 {
		return Percentage{}, ErrInvalidPercentage
	}
	return Percentage{value: value}, nil
}

// MustNewPercentage creates a Percentage or panics
func MustNewPercentage(value float64) Percentage {
	p, err := NewPercentage(value)
	if err != nil {
		panic(err)
	}
	return p
}

// Value returns the percentage value (0-100)
func (p Percentage) Value() float64 {
	return p.value
}

// Decimal returns the decimal representation (0-1)
func (p Percentage) Decimal() float64 {
	return p.value / 100
}

// Apply applies the percentage to a money amount
func (p Percentage) Apply(m Money) (Money, error) {
	return m.Percentage(p.value)
}

// Add adds two percentages
func (p Percentage) Add(other Percentage) (Percentage, error) {
	return NewPercentage(p.value + other.value)
}

// Subtract subtracts a percentage from this one
func (p Percentage) Subtract(other Percentage) (Percentage, error) {
	return NewPercentage(p.value - other.value)
}

// IsZero returns true if the percentage is zero
func (p Percentage) IsZero() bool {
	return p.value == 0
}

// Equals checks equality
func (p Percentage) Equals(other Percentage) bool {
	return p.value == other.value
}

// String returns a formatted string representation
func (p Percentage) String() string {
	return fmt.Sprintf("%.2f%%", p.value)
}

// BasisPoints represents a basis point value (1/100 of a percentage)
type BasisPoints struct {
	value int
}

// NewBasisPoints creates a new BasisPoints value
func NewBasisPoints(value int) (BasisPoints, error) {
	if value < 0 || value > 10000 {
		return BasisPoints{}, errors.New("basis points must be between 0 and 10000")
	}
	return BasisPoints{value: value}, nil
}

// Value returns the basis points value
func (bp BasisPoints) Value() int {
	return bp.value
}

// ToPercentage converts basis points to percentage
func (bp BasisPoints) ToPercentage() Percentage {
	return Percentage{value: float64(bp.value) / 100}
}

// Decimal returns the decimal representation
func (bp BasisPoints) Decimal() float64 {
	return float64(bp.value) / 10000
}

// String returns a formatted string
func (bp BasisPoints) String() string {
	return fmt.Sprintf("%d bps", bp.value)
}
