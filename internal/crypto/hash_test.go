package crypto

import (
	"testing"
)

func TestHash(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		wantHash [32]byte
	}{
		{
			name:  "empty input",
			input: []byte{},
			wantHash: [32]byte{
				0x5d, 0xf6, 0xe0, 0xe2, 0x76, 0x13, 0x59, 0xd3,
				0x0a, 0x82, 0x75, 0x05, 0x8e, 0x29, 0x9f, 0xcc,
				0x03, 0x81, 0x53, 0x45, 0x45, 0xf5, 0x5c, 0xf4,
				0x3e, 0x41, 0x98, 0x3f, 0x5d, 0x4c, 0x94, 0x56,
			},
		},
		{
			name:  "simple string",
			input: []byte("test"),
			wantHash: [32]byte{
				0x95, 0x4d, 0x5a, 0x49, 0xfd, 0x70, 0xd9, 0xb8,
				0xbc, 0xdb, 0x35, 0xd2, 0x52, 0x26, 0x78, 0x29,
				0x95, 0x7f, 0x7e, 0xf7, 0xfa, 0x6c, 0x74, 0xf8,
				0x84, 0x19, 0xbd, 0xc5, 0xe8, 0x22, 0x09, 0xf4,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Hash(tt.input)
			if got != tt.wantHash {
				t.Errorf("Hash() = %x, want %x", got, tt.wantHash)
			}
		})
	}
}

func TestKeyPairGeneration(t *testing.T) {
	privKey, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}
	if privKey == nil {
		t.Fatal("Generated private key is nil")
	}
}

func TestSignAndVerify(t *testing.T) {
	// Generate a key pair
	privKey, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	data := []byte("test message")

	// Sign
	signature, err := Sign(privKey, data)
	if err != nil {
		t.Fatalf("Failed to sign: %v", err)
	}

	// Verify
	if !Verify(&privKey.PublicKey, data, signature) {
		t.Error("Signature verification failed")
	}

	// Verify with wrong data
	wrongData := []byte("wrong message")
	if Verify(&privKey.PublicKey, wrongData, signature) {
		t.Error("Verification should fail with wrong data")
	}
}
