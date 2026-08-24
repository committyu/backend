package domain

import "testing"

func TestNewIDsGenerateValidDistinctUUIDs(t *testing.T) {
	userID := NewUserID()
	characterID := NewCharacterID()

	if !IsValidUserID(userID.String()) {
		t.Fatalf("NewUserID() returned an invalid UUID: %q", userID)
	}
	if !IsValidCharacterID(characterID.String()) {
		t.Fatalf("NewCharacterID() returned an invalid UUID: %q", characterID)
	}
	if userID.String() == characterID.String() {
		t.Fatal("generated user and character IDs must be distinct")
	}
}

func TestParseIDs(t *testing.T) {
	const input = "550E8400-E29B-41D4-A716-446655440000"
	const canonical = "550e8400-e29b-41d4-a716-446655440000"

	userID, err := ParseUserID(input)
	if err != nil {
		t.Fatalf("ParseUserID() error = %v", err)
	}
	if userID.String() != canonical {
		t.Fatalf("ParseUserID() = %q, want %q", userID, canonical)
	}

	characterID, err := ParseCharacterID(input)
	if err != nil {
		t.Fatalf("ParseCharacterID() error = %v", err)
	}
	if characterID.String() != canonical {
		t.Fatalf("ParseCharacterID() = %q, want %q", characterID, canonical)
	}
}

func TestInvalidIDs(t *testing.T) {
	const invalid = "not-a-uuid"

	if IsValidUserID(invalid) || IsValidCharacterID(invalid) {
		t.Fatal("invalid UUID was accepted")
	}
	if _, err := ParseUserID(invalid); err == nil {
		t.Fatal("ParseUserID() accepted an invalid UUID")
	}
	if _, err := ParseCharacterID(invalid); err == nil {
		t.Fatal("ParseCharacterID() accepted an invalid UUID")
	}
}