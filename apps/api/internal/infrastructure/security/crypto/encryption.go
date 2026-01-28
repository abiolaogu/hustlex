package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"io"

	"golang.org/x/crypto/argon2"
)

// EncryptionService provides field-level encryption for PII
type EncryptionService interface {
	Encrypt(plaintext []byte) ([]byte, error)
	Decrypt(ciphertext []byte) ([]byte, error)
	EncryptString(plaintext string) (string, error)
	DecryptString(ciphertext string) (string, error)
	DeriveKey(password string, salt []byte) []byte
}

// Argon2 parameters following OWASP recommendations
const (
	argon2Time    = 1
	argon2Memory  = 64 * 1024 // 64MB
	argon2Threads = 4
	argon2KeyLen  = 32
	saltLength    = 16
)

// AESEncryptionService implements EncryptionService using AES-256-GCM
type AESEncryptionService struct {
	key []byte
}

// NewAESEncryptionService creates a new AES encryption service
// Key must be 32 bytes for AES-256
func NewAESEncryptionService(key []byte) (*AESEncryptionService, error) {
	if len(key) != 32 {
		return nil, errors.New("key must be 32 bytes for AES-256")
	}
	return &AESEncryptionService{key: key}, nil
}

// NewAESEncryptionServiceFromPassword creates an encryption service from a password
func NewAESEncryptionServiceFromPassword(password string, salt []byte) (*AESEncryptionService, error) {
	if len(salt) < saltLength {
		return nil, errors.New("salt must be at least 16 bytes")
	}
	key := argon2.IDKey([]byte(password), salt, argon2Time, argon2Memory, argon2Threads, argon2KeyLen)
	return NewAESEncryptionService(key)
}

// Encrypt encrypts plaintext using AES-256-GCM
func (s *AESEncryptionService) Encrypt(plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(s.key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

// Decrypt decrypts ciphertext using AES-256-GCM
func (s *AESEncryptionService) Decrypt(ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(s.key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < gcm.NonceSize() {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:gcm.NonceSize()], ciphertext[gcm.NonceSize():]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

// EncryptString encrypts a string and returns base64-encoded ciphertext
func (s *AESEncryptionService) EncryptString(plaintext string) (string, error) {
	encrypted, err := s.Encrypt([]byte(plaintext))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

// DecryptString decrypts a base64-encoded ciphertext string
func (s *AESEncryptionService) DecryptString(ciphertext string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}
	decrypted, err := s.Decrypt(data)
	if err != nil {
		return "", err
	}
	return string(decrypted), nil
}

// DeriveKey uses Argon2id for key derivation (OWASP recommended)
func (s *AESEncryptionService) DeriveKey(password string, salt []byte) []byte {
	return argon2.IDKey([]byte(password), salt, argon2Time, argon2Memory, argon2Threads, argon2KeyLen)
}

// GenerateSalt generates a cryptographically secure random salt
func GenerateSalt() ([]byte, error) {
	salt := make([]byte, saltLength)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}
	return salt, nil
}

// GenerateKey generates a cryptographically secure random 256-bit key
func GenerateKey() ([]byte, error) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return nil, err
	}
	return key, nil
}

// HashPassword creates a secure password hash using Argon2id
// Returns base64-encoded salt+hash
func HashPassword(password string) (string, error) {
	salt := make([]byte, saltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, argon2Time, argon2Memory, argon2Threads, argon2KeyLen)

	// Combine salt and hash: salt (16 bytes) + hash (32 bytes)
	result := make([]byte, len(salt)+len(hash))
	copy(result, salt)
	copy(result[len(salt):], hash)

	return base64.StdEncoding.EncodeToString(result), nil
}

// VerifyPassword verifies a password against a hash using constant-time comparison
func VerifyPassword(password, encodedHash string) bool {
	data, err := base64.StdEncoding.DecodeString(encodedHash)
	if err != nil {
		return false
	}

	if len(data) < saltLength+argon2KeyLen {
		return false
	}

	salt := data[:saltLength]
	expectedHash := data[saltLength:]

	actualHash := argon2.IDKey([]byte(password), salt, argon2Time, argon2Memory, argon2Threads, argon2KeyLen)

	return subtle.ConstantTimeCompare(expectedHash, actualHash) == 1
}

// HashPIN creates a secure hash for transaction PINs
func HashPIN(pin string) (string, error) {
	return HashPassword(pin)
}

// VerifyPIN verifies a PIN against a hash
func VerifyPIN(pin, encodedHash string) bool {
	return VerifyPassword(pin, encodedHash)
}

// MaskPII masks personally identifiable information for logging
func MaskPII(value string, showLast int) string {
	if len(value) <= showLast {
		return "****"
	}
	masked := make([]byte, len(value))
	for i := 0; i < len(value)-showLast; i++ {
		masked[i] = '*'
	}
	copy(masked[len(value)-showLast:], value[len(value)-showLast:])
	return string(masked)
}

// MaskEmail masks an email address for logging
func MaskEmail(email string) string {
	for i, c := range email {
		if c == '@' {
			if i <= 2 {
				return "***" + email[i:]
			}
			return email[:2] + "***" + email[i:]
		}
	}
	return MaskPII(email, 4)
}

// MaskPhone masks a phone number for logging
func MaskPhone(phone string) string {
	return MaskPII(phone, 4)
}

// MaskBVN masks a BVN for logging
func MaskBVN(bvn string) string {
	return MaskPII(bvn, 4)
}

// MaskAccountNumber masks a bank account number for logging
func MaskAccountNumber(accountNumber string) string {
	return MaskPII(accountNumber, 4)
}
