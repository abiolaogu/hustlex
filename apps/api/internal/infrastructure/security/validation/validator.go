package validation

import (
	"errors"
	"fmt"
	"net"
	"net/mail"
	"regexp"
	"strings"
	"unicode"
)

// ValidationError holds validation errors for multiple fields
type ValidationError struct {
	Errors map[string]string `json:"errors"`
}

func (e *ValidationError) Error() string {
	if len(e.Errors) == 0 {
		return "validation failed"
	}
	var msgs []string
	for field, msg := range e.Errors {
		msgs = append(msgs, field+": "+msg)
	}
	return strings.Join(msgs, "; ")
}

func (e *ValidationError) HasErrors() bool {
	return len(e.Errors) > 0
}

func (e *ValidationError) Add(field, message string) {
	if e.Errors == nil {
		e.Errors = make(map[string]string)
	}
	e.Errors[field] = message
}

// NewValidationError creates a new validation error
func NewValidationError() *ValidationError {
	return &ValidationError{
		Errors: make(map[string]string),
	}
}

// Common validation patterns
var (
	// Nigerian phone number: +234XXXXXXXXXX or 0XXXXXXXXXX
	PhoneRegex = regexp.MustCompile(`^(\+234|0)[789][01]\d{8}$`)

	// Email validation (deprecated - use RFC-compliant mail.ParseAddress instead)
	// Kept for backward compatibility but not used in validation
	EmailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	// BVN: 11 digits
	BVNRegex = regexp.MustCompile(`^\d{11}$`)

	// NIN: 11 digits
	NINRegex = regexp.MustCompile(`^\d{11}$`)

	// Account number: 10 digits (NUBAN)
	AccountNumberRegex = regexp.MustCompile(`^\d{10}$`)

	// Bank code: 3 digits
	BankCodeRegex = regexp.MustCompile(`^\d{3}$`)

	// UUID v4
	UUIDRegex = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-4[0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$`)

	// Currency code: 3 uppercase letters
	CurrencyRegex = regexp.MustCompile(`^[A-Z]{3}$`)

	// Alphanumeric with underscores/hyphens (for references)
	ReferenceRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

	// SQL injection patterns to block
	SQLInjectionPatterns = []*regexp.Regexp{
		regexp.MustCompile(`(?i)(union\s+select|select\s+\*|drop\s+table|insert\s+into|delete\s+from|update\s+\w+\s+set)`),
		regexp.MustCompile(`(?i)(--|;|\/\*|\*\/|xp_|sp_)`),
	}

	// XSS patterns to block
	XSSPatterns = []*regexp.Regexp{
		regexp.MustCompile(`(?i)<script[^>]*>.*?</script>`),
		regexp.MustCompile(`(?i)javascript:`),
		regexp.MustCompile(`(?i)on\w+\s*=`),
		regexp.MustCompile(`(?i)<iframe`),
	}
)

// Validator provides input validation functions
type Validator struct {
	errors *ValidationError
}

// NewValidator creates a new validator instance
func NewValidator() *Validator {
	return &Validator{
		errors: NewValidationError(),
	}
}

// Errors returns the validation errors
func (v *Validator) Errors() *ValidationError {
	return v.errors
}

// HasErrors returns true if there are validation errors
func (v *Validator) HasErrors() bool {
	return v.errors.HasErrors()
}

// Required validates that a string field is not empty
func (v *Validator) Required(field, value string) *Validator {
	if strings.TrimSpace(value) == "" {
		v.errors.Add(field, "is required")
	}
	return v
}

// MinLength validates minimum string length
func (v *Validator) MinLength(field, value string, min int) *Validator {
	if len(value) < min {
		v.errors.Add(field, "must be at least "+itoa(min)+" characters")
	}
	return v
}

// MaxLength validates maximum string length
func (v *Validator) MaxLength(field, value string, max int) *Validator {
	if len(value) > max {
		v.errors.Add(field, "must be at most "+itoa(max)+" characters")
	}
	return v
}

// Length validates exact string length
func (v *Validator) Length(field, value string, length int) *Validator {
	if len(value) != length {
		v.errors.Add(field, "must be exactly "+itoa(length)+" characters")
	}
	return v
}

// Min validates minimum numeric value
func (v *Validator) Min(field string, value, min int64) *Validator {
	if value < min {
		v.errors.Add(field, "must be at least "+itoa64(min))
	}
	return v
}

// Max validates maximum numeric value
func (v *Validator) Max(field string, value, max int64) *Validator {
	if value > max {
		v.errors.Add(field, "must be at most "+itoa64(max))
	}
	return v
}

// Positive validates that a number is positive
func (v *Validator) Positive(field string, value int64) *Validator {
	if value <= 0 {
		v.errors.Add(field, "must be positive")
	}
	return v
}

// NonNegative validates that a number is non-negative
func (v *Validator) NonNegative(field string, value int64) *Validator {
	if value < 0 {
		v.errors.Add(field, "must not be negative")
	}
	return v
}

// Email validates email format using RFC 5321 compliant parser
func (v *Validator) Email(field, value string) *Validator {
	if value == "" {
		return v
	}

	// Use Go's standard library for RFC-compliant email parsing
	addr, err := mail.ParseAddress(value)
	if err != nil {
		v.errors.Add(field, "must be a valid email address")
		return v
	}

	// Additional RFC checks
	if len(addr.Address) > 254 {
		v.errors.Add(field, "email address too long (max 254 characters)")
		return v
	}

	// Validate domain part has at least one dot
	parts := strings.Split(addr.Address, "@")
	if len(parts) != 2 || !strings.Contains(parts[1], ".") {
		v.errors.Add(field, "must be a valid email address")
		return v
	}

	return v
}

// Phone validates Nigerian phone number format
func (v *Validator) Phone(field, value string) *Validator {
	if value != "" && !PhoneRegex.MatchString(value) {
		v.errors.Add(field, "must be a valid Nigerian phone number")
	}
	return v
}

// BVN validates Bank Verification Number format
func (v *Validator) BVN(field, value string) *Validator {
	if value != "" && !BVNRegex.MatchString(value) {
		v.errors.Add(field, "must be a valid 11-digit BVN")
	}
	return v
}

// NIN validates National Identification Number format
func (v *Validator) NIN(field, value string) *Validator {
	if value != "" && !NINRegex.MatchString(value) {
		v.errors.Add(field, "must be a valid 11-digit NIN")
	}
	return v
}

// AccountNumber validates NUBAN account number format
func (v *Validator) AccountNumber(field, value string) *Validator {
	if value != "" && !AccountNumberRegex.MatchString(value) {
		v.errors.Add(field, "must be a valid 10-digit account number")
	}
	return v
}

// BankCode validates bank code format
func (v *Validator) BankCode(field, value string) *Validator {
	if value != "" && !BankCodeRegex.MatchString(value) {
		v.errors.Add(field, "must be a valid 3-digit bank code")
	}
	return v
}

// UUID validates UUID v4 format
func (v *Validator) UUID(field, value string) *Validator {
	if value != "" && !UUIDRegex.MatchString(value) {
		v.errors.Add(field, "must be a valid UUID")
	}
	return v
}

// Currency validates currency code format
func (v *Validator) Currency(field, value string) *Validator {
	if value != "" && !CurrencyRegex.MatchString(value) {
		v.errors.Add(field, "must be a valid 3-letter currency code")
	}
	return v
}

// Reference validates reference format
func (v *Validator) Reference(field, value string) *Validator {
	if value != "" && !ReferenceRegex.MatchString(value) {
		v.errors.Add(field, "must contain only alphanumeric characters, underscores, and hyphens")
	}
	return v
}

// PIN validates PIN format (4-6 digits)
func (v *Validator) PIN(field, value string) *Validator {
	if value == "" {
		return v
	}
	if len(value) < 4 || len(value) > 6 {
		v.errors.Add(field, "must be 4-6 digits")
		return v
	}
	for _, c := range value {
		if !unicode.IsDigit(c) {
			v.errors.Add(field, "must contain only digits")
			return v
		}
	}
	return v
}

// Password validates password strength
func (v *Validator) Password(field, value string) *Validator {
	if value == "" {
		return v
	}

	var hasUpper, hasLower, hasDigit, hasSpecial bool
	for _, c := range value {
		switch {
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsDigit(c):
			hasDigit = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			hasSpecial = true
		}
	}

	if len(value) < 12 {
		v.errors.Add(field, "must be at least 12 characters")
	} else if !hasUpper || !hasLower || !hasDigit {
		v.errors.Add(field, "must contain uppercase, lowercase, and digits")
	} else if !hasSpecial {
		// Warn but don't fail
	}

	return v
}

// OneOf validates that value is one of allowed values
func (v *Validator) OneOf(field, value string, allowed []string) *Validator {
	if value == "" {
		return v
	}
	for _, a := range allowed {
		if value == a {
			return v
		}
	}
	v.errors.Add(field, "must be one of: "+strings.Join(allowed, ", "))
	return v
}

// NoSQLInjection validates that value doesn't contain SQL injection patterns
func (v *Validator) NoSQLInjection(field, value string) *Validator {
	for _, pattern := range SQLInjectionPatterns {
		if pattern.MatchString(value) {
			v.errors.Add(field, "contains invalid characters")
			return v
		}
	}
	return v
}

// NoXSS validates that value doesn't contain XSS patterns
func (v *Validator) NoXSS(field, value string) *Validator {
	for _, pattern := range XSSPatterns {
		if pattern.MatchString(value) {
			v.errors.Add(field, "contains invalid characters")
			return v
		}
	}
	return v
}

// SafeString validates that a string is safe (no SQL injection or XSS)
func (v *Validator) SafeString(field, value string) *Validator {
	return v.NoSQLInjection(field, value).NoXSS(field, value)
}

// Custom adds a custom validation
func (v *Validator) Custom(field string, valid bool, message string) *Validator {
	if !valid {
		v.errors.Add(field, message)
	}
	return v
}

// Validate returns error if validation failed
func (v *Validator) Validate() error {
	if v.errors.HasErrors() {
		return v.errors
	}
	return nil
}

// Helper functions

func itoa(i int) string {
	return itoa64(int64(i))
}

func itoa64(i int64) string {
	if i < 0 {
		return "-" + itoa64(-i)
	}
	if i < 10 {
		return string(rune('0' + i))
	}
	return itoa64(i/10) + string(rune('0'+i%10))
}

// Standalone validation functions

// ValidatePhone validates a Nigerian phone number
func ValidatePhone(phone string) error {
	if !PhoneRegex.MatchString(phone) {
		return errors.New("invalid phone number format")
	}
	return nil
}

// ValidateEmail validates an email address using RFC 5321 compliant parser
func ValidateEmail(email string) error {
	if email == "" {
		return errors.New("email is required")
	}

	// Use Go's standard library for RFC-compliant email parsing
	addr, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("invalid email format: %w", err)
	}

	// Additional RFC checks
	if len(addr.Address) > 254 {
		return errors.New("email address too long (max 254 characters)")
	}

	// Validate domain part has at least one dot
	parts := strings.Split(addr.Address, "@")
	if len(parts) != 2 || !strings.Contains(parts[1], ".") {
		return errors.New("invalid email format: domain must contain at least one dot")
	}

	return nil
}

// ValidateEmailWithDNS validates an email address and optionally checks DNS MX records
func ValidateEmailWithDNS(email string, checkDNS bool) error {
	// First do basic RFC validation
	if err := ValidateEmail(email); err != nil {
		return err
	}

	// Optional DNS MX record validation
	if checkDNS {
		addr, _ := mail.ParseAddress(email)
		parts := strings.Split(addr.Address, "@")
		domain := parts[1]

		mxRecords, err := net.LookupMX(domain)
		if err != nil || len(mxRecords) == 0 {
			return fmt.Errorf("email domain '%s' has no valid MX records", domain)
		}
	}

	return nil
}

// ValidateBVN validates a BVN
func ValidateBVN(bvn string) error {
	if !BVNRegex.MatchString(bvn) {
		return errors.New("invalid BVN format")
	}
	return nil
}

// ValidateAccountNumber validates a NUBAN account number
func ValidateAccountNumber(accountNumber string) error {
	if !AccountNumberRegex.MatchString(accountNumber) {
		return errors.New("invalid account number format")
	}
	return nil
}

// ValidateAmount validates a transaction amount
func ValidateAmount(amount int64, min, max int64) error {
	if amount < min {
		return errors.New("amount below minimum")
	}
	if amount > max {
		return errors.New("amount exceeds maximum")
	}
	return nil
}

// NormalizePhone normalizes a Nigerian phone number to +234 format
func NormalizePhone(phone string) string {
	phone = strings.TrimSpace(phone)
	if strings.HasPrefix(phone, "0") {
		return "+234" + phone[1:]
	}
	if !strings.HasPrefix(phone, "+") {
		return "+234" + phone
	}
	return phone
}

// SanitizeString removes potentially dangerous characters
func SanitizeString(s string) string {
	// Remove null bytes
	s = strings.ReplaceAll(s, "\x00", "")
	// Trim whitespace
	s = strings.TrimSpace(s)
	return s
}
