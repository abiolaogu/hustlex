package crypto

import (
	"testing"
)

func TestAESEncryptionService_EncryptDecrypt(t *testing.T) {
	key, err := GenerateKey()
	if err != nil {
		t.Fatalf("GenerateKey() error: %v", err)
	}

	svc, err := NewAESEncryptionService(key)
	if err != nil {
		t.Fatalf("NewAESEncryptionService() error: %v", err)
	}

	plaintext := []byte("sensitive data to encrypt")

	encrypted, err := svc.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encrypt() error: %v", err)
	}

	if string(encrypted) == string(plaintext) {
		t.Error("Encrypt() should not return plaintext")
	}

	decrypted, err := svc.Decrypt(encrypted)
	if err != nil {
		t.Fatalf("Decrypt() error: %v", err)
	}

	if string(decrypted) != string(plaintext) {
		t.Errorf("Decrypt() = %s, want %s", decrypted, plaintext)
	}
}

func TestAESEncryptionService_EncryptDecryptString(t *testing.T) {
	key, _ := GenerateKey()
	svc, _ := NewAESEncryptionService(key)

	plaintext := "user@example.com"

	encrypted, err := svc.EncryptString(plaintext)
	if err != nil {
		t.Fatalf("EncryptString() error: %v", err)
	}

	if encrypted == plaintext {
		t.Error("EncryptString() should not return plaintext")
	}

	decrypted, err := svc.DecryptString(encrypted)
	if err != nil {
		t.Fatalf("DecryptString() error: %v", err)
	}

	if decrypted != plaintext {
		t.Errorf("DecryptString() = %s, want %s", decrypted, plaintext)
	}
}

func TestAESEncryptionService_InvalidKeyLength(t *testing.T) {
	shortKey := []byte("too short")

	_, err := NewAESEncryptionService(shortKey)
	if err == nil {
		t.Error("NewAESEncryptionService() should error on invalid key length")
	}
}

func TestAESEncryptionService_DecryptInvalidData(t *testing.T) {
	key, _ := GenerateKey()
	svc, _ := NewAESEncryptionService(key)

	_, err := svc.Decrypt([]byte("invalid"))
	if err == nil {
		t.Error("Decrypt() should error on invalid ciphertext")
	}
}

func TestHashPassword_VerifyPassword(t *testing.T) {
	password := "SecureP@ssw0rd123!"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() error: %v", err)
	}

	if hash == password {
		t.Error("HashPassword() should not return plaintext")
	}

	// Verify correct password
	if !VerifyPassword(password, hash) {
		t.Error("VerifyPassword() should return true for correct password")
	}

	// Verify incorrect password
	if VerifyPassword("wrongpassword", hash) {
		t.Error("VerifyPassword() should return false for incorrect password")
	}
}

func TestHashPassword_UniqueHashes(t *testing.T) {
	password := "SamePassword123!"

	hash1, _ := HashPassword(password)
	hash2, _ := HashPassword(password)

	// Same password should produce different hashes (due to random salt)
	if hash1 == hash2 {
		t.Error("HashPassword() should produce unique hashes for same password")
	}

	// But both should verify correctly
	if !VerifyPassword(password, hash1) || !VerifyPassword(password, hash2) {
		t.Error("Both hashes should verify correctly")
	}
}

func TestVerifyPassword_InvalidHash(t *testing.T) {
	if VerifyPassword("password", "invalid-base64!!!") {
		t.Error("VerifyPassword() should return false for invalid hash")
	}

	if VerifyPassword("password", "c2hvcnQ=") { // "short" in base64
		t.Error("VerifyPassword() should return false for truncated hash")
	}
}

func TestHashPIN_VerifyPIN(t *testing.T) {
	pin := "1234"

	hash, err := HashPIN(pin)
	if err != nil {
		t.Fatalf("HashPIN() error: %v", err)
	}

	if !VerifyPIN(pin, hash) {
		t.Error("VerifyPIN() should return true for correct PIN")
	}

	if VerifyPIN("0000", hash) {
		t.Error("VerifyPIN() should return false for incorrect PIN")
	}
}

func TestGenerateSalt(t *testing.T) {
	salt1, err := GenerateSalt()
	if err != nil {
		t.Fatalf("GenerateSalt() error: %v", err)
	}

	if len(salt1) != saltLength {
		t.Errorf("GenerateSalt() length = %d, want %d", len(salt1), saltLength)
	}

	salt2, _ := GenerateSalt()
	if string(salt1) == string(salt2) {
		t.Error("GenerateSalt() should produce unique salts")
	}
}

func TestGenerateKey(t *testing.T) {
	key1, err := GenerateKey()
	if err != nil {
		t.Fatalf("GenerateKey() error: %v", err)
	}

	if len(key1) != 32 {
		t.Errorf("GenerateKey() length = %d, want 32", len(key1))
	}

	key2, _ := GenerateKey()
	if string(key1) == string(key2) {
		t.Error("GenerateKey() should produce unique keys")
	}
}

func TestMaskPII(t *testing.T) {
	tests := []struct {
		value    string
		showLast int
		want     string
	}{
		{"1234567890", 4, "******7890"},
		{"12345678901", 4, "*******8901"},
		{"123", 4, "****"},
		{"", 4, "****"},
	}

	for _, tt := range tests {
		got := MaskPII(tt.value, tt.showLast)
		if got != tt.want {
			t.Errorf("MaskPII(%q, %d) = %q, want %q", tt.value, tt.showLast, got, tt.want)
		}
	}
}

func TestMaskEmail(t *testing.T) {
	tests := []struct {
		email string
		want  string
	}{
		{"user@example.com", "us***@example.com"},
		{"a@b.com", "***@b.com"},
		{"ab@test.org", "***@test.org"},
		{"invalid", "***alid"},
	}

	for _, tt := range tests {
		got := MaskEmail(tt.email)
		if got != tt.want {
			t.Errorf("MaskEmail(%q) = %q, want %q", tt.email, got, tt.want)
		}
	}
}

func TestMaskPhone(t *testing.T) {
	got := MaskPhone("+2348012345678")
	if got != "**********5678" {
		t.Errorf("MaskPhone() = %q, want **********5678", got)
	}
}

func TestMaskBVN(t *testing.T) {
	got := MaskBVN("12345678901")
	if got != "*******8901" {
		t.Errorf("MaskBVN() = %q, want *******8901", got)
	}
}

func TestMaskAccountNumber(t *testing.T) {
	got := MaskAccountNumber("0123456789")
	if got != "******6789" {
		t.Errorf("MaskAccountNumber() = %q, want ******6789", got)
	}
}

func TestNewAESEncryptionServiceFromPassword(t *testing.T) {
	password := "my-secret-password"
	salt, _ := GenerateSalt()

	svc, err := NewAESEncryptionServiceFromPassword(password, salt)
	if err != nil {
		t.Fatalf("NewAESEncryptionServiceFromPassword() error: %v", err)
	}

	plaintext := "test data"
	encrypted, _ := svc.EncryptString(plaintext)

	// Create another service with same password and salt
	svc2, _ := NewAESEncryptionServiceFromPassword(password, salt)
	decrypted, err := svc2.DecryptString(encrypted)
	if err != nil {
		t.Fatalf("DecryptString() error: %v", err)
	}

	if decrypted != plaintext {
		t.Errorf("DecryptString() = %s, want %s", decrypted, plaintext)
	}
}

func TestNewAESEncryptionServiceFromPassword_InvalidSalt(t *testing.T) {
	_, err := NewAESEncryptionServiceFromPassword("password", []byte("short"))
	if err == nil {
		t.Error("NewAESEncryptionServiceFromPassword() should error on short salt")
	}
}
