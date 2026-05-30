package crypto

import "testing"

func TestTokenCipherEncryptDecrypt(t *testing.T) {
	cipher, err := NewTokenCipher("dev-secret")
	if err != nil {
		t.Fatalf("NewTokenCipher() error = %v", err)
	}

	encrypted, err := cipher.Encrypt("github-token")
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}

	if encrypted == "github-token" {
		t.Fatal("Encrypt() returned plaintext")
	}

	decrypted, err := cipher.Decrypt(encrypted)
	if err != nil {
		t.Fatalf("Decrypt() error = %v", err)
	}

	if decrypted != "github-token" {
		t.Fatalf("Decrypt() = %q, want %q", decrypted, "github-token")
	}
}

func TestTokenCipherDecryptWithDifferentKeyFails(t *testing.T) {
	cipher, err := NewTokenCipher("dev-secret")
	if err != nil {
		t.Fatalf("NewTokenCipher() error = %v", err)
	}

	encrypted, err := cipher.Encrypt("github-token")
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}

	otherCipher, err := NewTokenCipher("another-secret")
	if err != nil {
		t.Fatalf("NewTokenCipher() error = %v", err)
	}

	if _, err := otherCipher.Decrypt(encrypted); err == nil {
		t.Fatal("Decrypt() error = nil, want error")
	}
}

func TestTokenCipherEmptyToken(t *testing.T) {
	cipher, err := NewTokenCipher("dev-secret")
	if err != nil {
		t.Fatalf("NewTokenCipher() error = %v", err)
	}

	encrypted, err := cipher.Encrypt("")
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}
	if encrypted != "" {
		t.Fatalf("Encrypt() = %q, want empty string", encrypted)
	}

	decrypted, err := cipher.Decrypt("")
	if err != nil {
		t.Fatalf("Decrypt() error = %v", err)
	}
	if decrypted != "" {
		t.Fatalf("Decrypt() = %q, want empty string", decrypted)
	}
}
