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
		{"valid with hyphen", "user-name@example.com", false},

		// Invalid emails
		{"invalid no @", "testexample.com", true},
		{"invalid no domain", "test@", true},
		{"invalid no local", "@example.com", true},
		{"invalid double @", "test@@example.com", true},
		{"invalid spaces", "test @example.com", true},
		{"invalid missing TLD", "test@example", true},
		{"too long email", "a" + "verylongemailaddressthatexceedsthemaximumlengthof254charactersallowedbyRFC5321standardforvalidemailaddressesthisisjustatesttomakesureweproperlyvali" + "dateemaillengthaccordingtothestandardandrejectemailsthataretoolongwhichcouldcauseproblemswithemailserversanddelivery@example.com", true},

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
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		// Valid RFC 5321 compliant emails
		{"simple valid", "test@example.com", false},
		{"with subdomain", "test@mail.example.com", false},
		{"with plus addressing", "user+tag@example.com", false},
		{"with dots", "first.last@example.com", false},
		{"with numbers", "user123@example456.com", false},
		{"with hyphens", "user-name@example-domain.com", false},
		{"with underscores", "user_name@example.com", false},

		// Invalid formats
		{"no @", "testexample.com", true},
		{"no domain", "test@", true},
		{"no local part", "@example.com", true},
		{"double @", "test@@example.com", true},
		{"spaces", "test @example.com", true},
		{"missing TLD", "test@example", true},
		{"empty string", "", true},

		// Length validation
		{"local part too long", "a123456789012345678901234567890123456789012345678901234567890123456@example.com", true},
		{"total too long", "test@" + "verylongdomainname123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789.com", true},

		// Edge cases from RFC
		{"quoted local part", `"test user"@example.com`, false},
		{"IP address domain", "test@[192.168.1.1]", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEmail(tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateEmail(%q) error = %v, wantErr %v", tt.email, err, tt.wantErr)
			}
		})
	}
}

func TestValidateEmailWithDNS(t *testing.T) {
	// Test without DNS check (should always pass for valid format)
	err := ValidateEmailWithDNS("test@example.com", false)
	if err != nil {
		t.Errorf("ValidateEmailWithDNS() without DNS check failed: %v", err)
	}

	// Test with DNS check for a known domain (may fail in CI without network)
	// This test is skipped in environments without network access
	t.Run("with DNS check", func(t *testing.T) {
		t.Skip("Skipping DNS check test - requires network access")
		// err := ValidateEmailWithDNS("test@gmail.com", true)
		// This should pass since gmail.com has valid MX records
	})

	// Test invalid format with DNS check (should fail on format before DNS)
	err = ValidateEmailWithDNS("invalid", true)
	if err == nil {
		t.Error("ValidateEmailWithDNS() expected error for invalid format")
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
