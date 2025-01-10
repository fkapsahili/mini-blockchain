package crypto

import (
	"crypto/ecdh"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"
)

// Hash performs a double SHA256 hash of the input data
func Hash(data []byte) [32]byte {
	firstHash := sha256.Sum256(data)
	return sha256.Sum256(firstHash[:])
}

// GenerateKeyPair creates a new key pair for signing
func GenerateKeyPair() (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
}

// Sign creates a digital signature of the data
func Sign(privateKey *ecdsa.PrivateKey, data []byte) ([]byte, error) {
	return ecdsa.SignASN1(rand.Reader, privateKey, data)
}

// Verify checks if the signature is valid for the data
func Verify(publicKey *ecdsa.PublicKey, data []byte, signature []byte) bool {
	return ecdsa.VerifyASN1(publicKey, data, signature)
}

// HashToAddress converts a public key to a blockchain address
func HashToAddress(publicKey *ecdsa.PublicKey) []byte {
	// Convert ECDSA public key to ECDH public key
	curve := ecdh.P256()

	pubBytes := make([]byte, 0, 64)
	pubBytes = append(pubBytes, publicKey.X.Bytes()...)
	pubBytes = append(pubBytes, publicKey.Y.Bytes()...)

	ecdhPub, err := curve.NewPublicKey(pubBytes)
	if err != nil {
		return nil
	}

	pubKeyBytes := ecdhPub.Bytes()

	// SHA256 hash of the public key
	sha256Hash := sha256.Sum256(pubKeyBytes)

	address := sha256Hash[:20]

	return address
}

// CalculateMerkleRoot computes the merkle root of a slice of hashes
func CalculateMerkleRoot(hashes [][32]byte) [32]byte {
	if len(hashes) == 0 {
		return sha256.Sum256([]byte{})
	}

	current := make([][32]byte, len(hashes))
	copy(current, hashes)

	// Keep combining pairs until we have one hash
	for len(current) > 1 {
		if len(current)%2 != 0 {
			// Duplicate last hash if odd
			current = append(current, current[len(current)-1])
		}

		next := make([][32]byte, len(current)/2)
		for i := 0; i < len(current); i += 2 {
			combined := append(current[i][:], current[i+1][:]...)
			next[i/2] = Hash(combined)
		}
		current = next
	}

	return current[0]
}

// CheckProofOfWork verifies if a hash meets the difficulty target
func CheckProofOfWork(hash [32]byte, difficulty uint32) bool {
	// Convert difficulty to target: more zeros = more difficult
	target := byte(256 - difficulty)

	return hash[0] < target
}

// GenerateProofOfWork finds a nonce that makes the block hash meet difficulty
func GenerateProofOfWork(data []byte, difficulty uint32) (uint32, [32]byte) {
	var nonce uint32
	var hash [32]byte

	for {
		// Append nonce to data
		withNonce := append(data, byte(nonce))
		hash = Hash(withNonce)

		if CheckProofOfWork(hash, difficulty) {
			return nonce, hash
		}

		nonce++
	}
}

// BytesToPublicKey converts a byte slice to an ECDSA public key
func BytesToPublicKey(pub []byte) (*ecdsa.PublicKey, error) {
	curve := ecdh.P256()
	ecdhPub, err := curve.NewPublicKey(pub)
	if err != nil {
		return nil, fmt.Errorf("invalid public key bytes: %w", err)
	}

	// Get raw bytes
	rawBytes := ecdhPub.Bytes()

	if len(rawBytes) != 65 || rawBytes[0] != 0x04 {
		return nil, fmt.Errorf("invalid public key format")
	}

	x := new(big.Int).SetBytes(rawBytes[1:33])
	y := new(big.Int).SetBytes(rawBytes[33:65])

	return &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}, nil
}
