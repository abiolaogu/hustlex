package validation

import (
	"testing"
)

func TestValidator_Required(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"valid value", "test", false},
		{"empty string", "", true},
		{"whitespace only", "   ", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewValidator()
			v.Required("field", tt.value)
			if tt.wantErr && !v.HasErrors() {
				t.Error("Expected error, got none")
			}
			if !tt.wantErr && v.HasErrors() {
				t.Errorf("Expected no error, got: %v", v.Errors())
			}
		})
	}
}

func TestValidator_MinLength(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		min     int
		wantErr bool
	}{
		{"meets minimum", "test", 4, false},
		{"exceeds minimum", "testing", 4, false},
		{"below minimum", "hi", 4, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewValidator()
			v.MinLength("field", tt.value, tt.min)
			if tt.wantErr != v.HasErrors() {
				t.Errorf("MinLength() error = %v, wantErr %v", v.HasErrors(), tt.wantErr)
			}
		})
	}
}

func TestValidator_MaxLength(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		max     int
		wantErr bool
	}{
		{"within limit", "test", 10, false},
		{"at limit", "test", 4, false},
		{"exceeds limit", "testing", 4, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewValidator()
			v.MaxLength("field", tt.value, tt.max)
			if tt.wantErr != v.HasErrors() {
				t.Errorf("MaxLength() error = %v, wantErr %v", v.HasErrors(), tt.wantErr)
			}
		})
	}
}

func TestValidator_Positive(t *testing.T) {
	tests := []struct {
		name    string
		value   int64
		wantErr bool
	}{
		{"positive", 100, false},
		{"zero", 0, true},
		{"negative", -100, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewValidator()
			v.Positive("amount", tt.value)
			if tt.wantErr != v.HasErrors() {
				t.Errorf("Positive() error = %v, wantErr %v", v.HasErrors(), tt.wantErr)
			}
		})
	}
}

func TestValidator_Email(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		// Valid emails
		{"valid email", "test@example.com", false},
		{"valid with subdomain", "test@sub.example.com", false},
		{"valid with plus", "user+tag@example.com", false},
		{"valid with dots", "first.last@example.com", false},
		{"valid with numbers", "user123@example.com", false},
		{"valid with hyphens", "user-name@ex-ample.com", false},
		{"valid with underscore", "user_name@example.com", false},
		{"valid long domain", "test@subdomain.example.co.uk", false},

		// Invalid emails
		{"invalid no @", "testexample.com", true},
		{"invalid no domain", "test@", true},
		{"invalid no local part", "@example.com", true},
		{"invalid double @", "test@@example.com", true},
		{"invalid no TLD", "test@example", true},
		{"invalid consecutive dots in domain", "test@exam..ple.com", true},
		{"invalid domain starts with dot", "test@.example.com", true},
		{"invalid domain ends with dot", "test@example.com.", true},
		{"invalid domain starts with hyphen", "test@-example.com", true},
		{"invalid spaces", "test @example.com", true},
		{"invalid special chars", "test!#$@example.com", true},

		// Edge cases
		{"empty (skip)", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewValidator()
			v.Email("email", tt.value)
			if tt.wantErr != v.HasErrors() {
				t.Errorf("Email(%q) error = %v, wantErr %v", tt.value, v.HasErrors(), tt.wantErr)
			}
		})
	}
}

func TestValidator_Phone(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"valid +234", "+2348012345678", false},
		{"valid 0", "08012345678", false},
		{"valid 070", "07012345678", false},
		{"valid 090", "09012345678", false},
		{"invalid too short", "0801234567", true},
		{"invalid too long", "080123456789", true},
		{"invalid prefix", "06012345678", true},
		{"empty (skip)", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewValidator()
			v.Phone("phone", tt.value)
			if tt.wantErr != v.HasErrors() {
				t.Errorf("Phone(%q) error = %v, wantErr %v", tt.value, v.HasErrors(), tt.wantErr)
			}
		})
	}
}

func TestValidator_BVN(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"valid 11 digits", "12345678901", false},
		{"invalid 10 digits", "1234567890", true},
		{"invalid 12 digits", "123456789012", true},
		{"invalid with letters", "1234567890a", true},
		{"empty (skip)", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewValidator()
			v.BVN("bvn", tt.value)
			if tt.wantErr != v.HasErrors() {
				t.Errorf("BVN(%q) error = %v, wantErr %v", tt.value, v.HasErrors(), tt.wantErr)
			}
		})
	}
}

func TestValidator_AccountNumber(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"valid 10 digits", "1234567890", false},
		{"invalid 9 digits", "123456789", true},
		{"invalid 11 digits", "12345678901", true},
		{"invalid with letters", "123456789a", true},
		{"empty (skip)", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewValidator()
			v.AccountNumber("account", tt.value)
			if tt.wantErr != v.HasErrors() {
				t.Errorf("AccountNumber(%q) error = %v, wantErr %v", tt.value, v.HasErrors(), tt.wantErr)
			}
		})
	}
}

func TestValidator_BankCode(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"valid 3 digits", "058", false},
		{"invalid 2 digits", "05", true},
		{"invalid 4 digits", "0580", true},
		{"empty (skip)", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewValidator()
			v.BankCode("bank_code", tt.value)
			if tt.wantErr != v.HasErrors() {
				t.Errorf("BankCode(%q) error = %v, wantErr %v", tt.value, v.HasErrors(), tt.wantErr)
			}
		})
	}
}

func TestValidator_PIN(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"valid 4 digits", "1234", false},
		{"valid 6 digits", "123456", false},
		{"invalid 3 digits", "123", true},
		{"invalid 7 digits", "1234567", true},
		{"invalid with letters", "12a4", true},
		{"empty (skip)", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewValidator()
			v.PIN("pin", tt.value)
			if tt.wantErr != v.HasErrors() {
				t.Errorf("PIN(%q) error = %v, wantErr %v", tt.value, v.HasErrors(), tt.wantErr)
			}
		})
	}
}

func TestValidator_Password(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"valid password", "SecurePass123!", false},
		{"valid long", "VerySecurePassword123", false},
		{"too short", "Short1!", true},
		{"no uppercase", "lowercase123!", true},
		{"no lowercase", "UPPERCASE123!", true},
		{"no digits", "NoDigitsHere!", true},
		{"empty (skip)", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewValidator()
			v.Password("password", tt.value)
			if tt.wantErr != v.HasErrors() {
				t.Errorf("Password(%q) error = %v, wantErr %v", tt.value, v.HasErrors(), tt.wantErr)
			}
		})
	}
}

func TestValidator_OneOf(t *testing.T) {
	allowed := []string{"NGN", "USD", "GBP"}

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"valid NGN", "NGN", false},
		{"valid USD", "USD", false},
		{"invalid EUR", "EUR", true},
		{"empty (skip)", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewValidator()
			v.OneOf("currency", tt.value, allowed)
			if tt.wantErr != v.HasErrors() {
				t.Errorf("OneOf(%q) error = %v, wantErr %v", tt.value, v.HasErrors(), tt.wantErr)
			}
		})
	}
}

func TestValidator_NoSQLInjection(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"safe string", "normal text", false},
		{"sql union", "'; UNION SELECT * FROM users--", true},
		{"sql drop", "'; DROP TABLE users;", true},
		{"sql comment", "admin'--", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewValidator()
			v.NoSQLInjection("input", tt.value)
			if tt.wantErr != v.HasErrors() {
				t.Errorf("NoSQLInjection(%q) error = %v, wantErr %v", tt.value, v.HasErrors(), tt.wantErr)
			}
		})
	}
}

func TestValidator_NoXSS(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"safe string", "normal text", false},
		{"script tag", "<script>alert('xss')</script>", true},
		{"javascript uri", "javascript:alert(1)", true},
		{"onclick handler", "onclick=alert(1)", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewValidator()
			v.NoXSS("input", tt.value)
			if tt.wantErr != v.HasErrors() {
				t.Errorf("NoXSS(%q) error = %v, wantErr %v", tt.value, v.HasErrors(), tt.wantErr)
			}
		})
	}
}

func TestValidator_Chain(t *testing.T) {
	v := NewValidator()
	v.Required("name", "John").
		MinLength("name", "John", 2).
		MaxLength("name", "John", 50)

	if v.HasErrors() {
		t.Errorf("Expected no errors, got: %v", v.Errors())
	}
}

func TestValidator_ChainWithErrors(t *testing.T) {
	v := NewValidator()
	v.Required("name", "").
		Positive("amount", -100).
		Email("email", "invalid")

	if !v.HasErrors() {
		t.Error("Expected errors, got none")
	}

	errors := v.Errors().Errors
	if _, ok := errors["name"]; !ok {
		t.Error("Expected error for 'name'")
	}
	if _, ok := errors["amount"]; !ok {
		t.Error("Expected error for 'amount'")
	}
	if _, ok := errors["email"]; !ok {
		t.Error("Expected error for 'email'")
	}
}

func TestValidatePhone(t *testing.T) {
	if err := ValidatePhone("+2348012345678"); err != nil {
		t.Errorf("ValidatePhone() unexpected error: %v", err)
	}
	if err := ValidatePhone("invalid"); err == nil {
		t.Error("ValidatePhone() expected error for invalid phone")
	}
}

func TestValidateEmail(t *testing.T) {
	if err := ValidateEmail("test@example.com"); err != nil {
		t.Errorf("ValidateEmail() unexpected error: %v", err)
	}
	if err := ValidateEmail("invalid"); err == nil {
		t.Error("ValidateEmail() expected error for invalid email")
	}
}

func TestValidateEmailRFC(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		// Valid RFC-compliant emails
		{"valid simple", "user@example.com", false},
		{"valid with subdomain", "user@mail.example.com", false},
		{"valid with plus", "user+tag@example.com", false},
		{"valid with dots", "first.last@example.com", false},
		{"valid with numbers", "user123@example.com", false},
		{"valid with hyphens", "user-name@example.com", false},
		{"valid with underscore", "user_name@example.com", false},
		{"valid international TLD", "user@example.co.uk", false},
		{"valid with name", "John Doe <john@example.com>", false},

		// Invalid emails
		{"invalid no @", "userexample.com", true},
		{"invalid double @", "user@@example.com", true},
		{"invalid no domain", "user@", true},
		{"invalid no local", "@example.com", true},
		{"invalid no TLD", "user@example", true},
		{"invalid consecutive dots", "user@exam..ple.com", true},
		{"invalid domain starts with dot", "user@.example.com", true},
		{"invalid domain ends with dot", "user@example.com.", true},
		{"invalid domain starts with hyphen", "user@-example.com", true},
		{"invalid domain ends with hyphen", "user@example-.com", true},
		{"invalid spaces in address", "user name@example.com", true},
		{"invalid empty", "", true},

		// RFC length limits
		{"valid max local length", strings.Repeat("a", 64) + "@example.com", false},
		{"invalid local too long", strings.Repeat("a", 65) + "@example.com", true},
		{"valid max domain length", "user@" + strings.Repeat("a", 240) + ".com", false},
		{"invalid domain too long", "user@" + strings.Repeat("a", 250) + ".com", true},
		{"valid max total length", strings.Repeat("a", 64) + "@" + strings.Repeat("b", 180) + ".com", false},
		{"invalid total too long", strings.Repeat("a", 64) + "@" + strings.Repeat("b", 186) + ".com", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEmailRFC(tt.email)
			if tt.wantErr && err == nil {
				t.Errorf("ValidateEmailRFC(%q) expected error, got nil", tt.email)
			}
			if !tt.wantErr && err != nil {
				t.Errorf("ValidateEmailRFC(%q) unexpected error: %v", tt.email, err)
			}
		})
	}
}

func TestValidateEmailWithDNS(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		checkDNS bool
		wantErr  bool
	}{
		// Without DNS check (should validate format only)
		{"valid format no DNS check", "user@example.com", false, false},
		{"invalid format no DNS check", "invalid", false, true},

		// With DNS check (these tests may fail if network is unavailable)
		{"valid real domain", "test@gmail.com", true, false},
		{"valid real domain 2", "test@yahoo.com", true, false},
		{"invalid fake domain", "test@thisisafakedomain12345.com", true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEmailWithDNS(tt.email, tt.checkDNS)
			if tt.wantErr && err == nil {
				t.Errorf("ValidateEmailWithDNS(%q, %v) expected error, got nil", tt.email, tt.checkDNS)
			}
			if !tt.wantErr && err != nil {
				// Skip DNS tests that fail due to network issues
				if tt.checkDNS && strings.Contains(err.Error(), "domain") {
					t.Skipf("Skipping DNS test due to network issue: %v", err)
				} else {
					t.Errorf("ValidateEmailWithDNS(%q, %v) unexpected error: %v", tt.email, tt.checkDNS, err)
				}
			}
		})
	}
}

func TestValidateAmount(t *testing.T) {
	if err := ValidateAmount(5000, 100, 1000000); err != nil {
		t.Errorf("ValidateAmount() unexpected error: %v", err)
	}
	if err := ValidateAmount(50, 100, 1000000); err == nil {
		t.Error("ValidateAmount() expected error for amount below minimum")
	}
	if err := ValidateAmount(2000000, 100, 1000000); err == nil {
		t.Error("ValidateAmount() expected error for amount above maximum")
	}
}

func TestNormalizePhone(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"08012345678", "+2348012345678"},
		{"+2348012345678", "+2348012345678"},
		{"8012345678", "+2348012345678"},
		{" 08012345678 ", "+2348012345678"},
	}

	for _, tt := range tests {
		result := NormalizePhone(tt.input)
		if result != tt.expected {
			t.Errorf("NormalizePhone(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestSanitizeString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"normal", "normal"},
		{"  trimmed  ", "trimmed"},
		{"with\x00null", "withnull"},
	}

	for _, tt := range tests {
		result := SanitizeString(tt.input)
		if result != tt.expected {
			t.Errorf("SanitizeString(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestValidationError_Error(t *testing.T) {
	ve := NewValidationError()
	ve.Add("field1", "error1")
	ve.Add("field2", "error2")

	errStr := ve.Error()
	if errStr == "" {
		t.Error("Expected non-empty error string")
	}
}

func TestValidationError_HasErrors(t *testing.T) {
	ve := NewValidationError()
	if ve.HasErrors() {
		t.Error("Expected no errors initially")
	}

	ve.Add("field", "error")
	if !ve.HasErrors() {
		t.Error("Expected errors after adding")
	}
}
