package valueobject

import (
	"errors"
	"regexp"
	"strings"
)

var (
	ErrInvalidPhoneNumber = errors.New("invalid phone number")
	ErrInvalidEmail       = errors.New("invalid email address")
)

// PhoneNumber represents a validated Nigerian phone number
type PhoneNumber struct {
	value string // Stored in E.164 format: +234XXXXXXXXXX
}

// Nigerian phone number patterns
var (
	nigerianPhonePattern = regexp.MustCompile(`^(\+234|234|0)?([789][01]\d{8})$`)
)

// NewPhoneNumber creates and validates a Nigerian phone number
func NewPhoneNumber(phone string) (PhoneNumber, error) {
	// Remove spaces and dashes
	cleaned := strings.ReplaceAll(phone, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "-", "")

	matches := nigerianPhonePattern.FindStringSubmatch(cleaned)
	if matches == nil {
		return PhoneNumber{}, ErrInvalidPhoneNumber
	}

	// Convert to E.164 format
	normalized := "+234" + matches[2]

	return PhoneNumber{value: normalized}, nil
}

// MustNewPhoneNumber creates a PhoneNumber or panics
func MustNewPhoneNumber(phone string) PhoneNumber {
	p, err := NewPhoneNumber(phone)
	if err != nil {
		panic(err)
	}
	return p
}

// String returns the E.164 format
func (p PhoneNumber) String() string {
	return p.value
}

// Local returns the local format (0XXXXXXXXXX)
func (p PhoneNumber) Local() string {
	if len(p.value) > 4 {
		return "0" + p.value[4:]
	}
	return p.value
}

// Masked returns a masked version for display
func (p PhoneNumber) Masked() string {
	if len(p.value) < 8 {
		return p.value
	}
	return p.value[:7] + "****" + p.value[len(p.value)-2:]
}

// IsEmpty returns true if the phone number is empty
func (p PhoneNumber) IsEmpty() bool {
	return p.value == ""
}

// Equals checks equality
func (p PhoneNumber) Equals(other PhoneNumber) bool {
	return p.value == other.value
}

// Email represents a validated email address
type Email struct {
	value string
}

var emailPattern = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// NewEmail creates and validates an email address
func NewEmail(email string) (Email, error) {
	cleaned := strings.TrimSpace(strings.ToLower(email))

	if !emailPattern.MatchString(cleaned) {
		return Email{}, ErrInvalidEmail
	}

	return Email{value: cleaned}, nil
}

// MustNewEmail creates an Email or panics
func MustNewEmail(email string) Email {
	e, err := NewEmail(email)
	if err != nil {
		panic(err)
	}
	return e
}

// String returns the email address
func (e Email) String() string {
	return e.value
}

// Domain returns the domain part of the email
func (e Email) Domain() string {
	parts := strings.Split(e.value, "@")
	if len(parts) == 2 {
		return parts[1]
	}
	return ""
}

// Masked returns a masked version for display
func (e Email) Masked() string {
	parts := strings.Split(e.value, "@")
	if len(parts) != 2 {
		return e.value
	}

	local := parts[0]
	domain := parts[1]

	if len(local) <= 2 {
		return local + "***@" + domain
	}

	return local[:2] + "***@" + domain
}

// IsEmpty returns true if the email is empty
func (e Email) IsEmpty() bool {
	return e.value == ""
}

// Equals checks equality
func (e Email) Equals(other Email) bool {
	return e.value == other.value
}

// FullName represents a person's full name
type FullName struct {
	value string
}

// NewFullName creates a validated full name
func NewFullName(name string) (FullName, error) {
	cleaned := strings.TrimSpace(name)

	if len(cleaned) < 2 {
		return FullName{}, errors.New("name must be at least 2 characters")
	}

	if len(cleaned) > 100 {
		return FullName{}, errors.New("name must be at most 100 characters")
	}

	return FullName{value: cleaned}, nil
}

// MustNewFullName creates a FullName or panics
func MustNewFullName(name string) FullName {
	fn, err := NewFullName(name)
	if err != nil {
		panic(err)
	}
	return fn
}

// String returns the full name
func (n FullName) String() string {
	return n.value
}

// FirstName returns the first word as first name
func (n FullName) FirstName() string {
	parts := strings.Fields(n.value)
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}

// LastName returns the last word as last name
func (n FullName) LastName() string {
	parts := strings.Fields(n.value)
	if len(parts) > 1 {
		return parts[len(parts)-1]
	}
	return ""
}

// Initials returns the initials
func (n FullName) Initials() string {
	parts := strings.Fields(n.value)
	var initials string
	for _, part := range parts {
		if len(part) > 0 {
			initials += strings.ToUpper(string(part[0]))
		}
	}
	return initials
}

// IsEmpty returns true if the name is empty
func (n FullName) IsEmpty() bool {
	return n.value == ""
}

// Equals checks equality
func (n FullName) Equals(other FullName) bool {
	return n.value == other.value
}
