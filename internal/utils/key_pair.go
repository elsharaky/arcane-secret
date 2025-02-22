package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"golang.org/x/crypto/ssh"
)

type KeyPair struct {
	PrivateKey  []byte
	PublicKey   []byte
	Fingerprint *string
}

func GenerateKeyPair(algorithm string, size *int, sshFormat bool) (*KeyPair, error) {
	switch algorithm {
	case "RSA":
		if size == nil {
			size = new(int)
			*size = 2048
		}

		return GenerateRSAKeyPair(*size, sshFormat)
	default:
		return nil, fmt.Errorf("unsupported algorithm: %s", algorithm)
	}
}

func GenerateRSAKeyPair(size int, sshFormat bool) (*KeyPair, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, size)
	if err != nil {
		return nil, err
	}

	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	var (
		publicKeyPEM []byte
		fingerprint  *string
	)

	if !sshFormat {
		publicKey := privateKey.Public()

		publicKeyPEM = pem.EncodeToMemory(&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: x509.MarshalPKCS1PublicKey(publicKey.(*rsa.PublicKey)),
		})
	} else {
		publicKey, err := ssh.NewPublicKey(&privateKey.PublicKey)
		if err != nil {
			return nil, err
		}

		publicKeyPEM = ssh.MarshalAuthorizedKey(publicKey)
		fingerprint = new(string)
		*fingerprint = ssh.FingerprintSHA256(publicKey)
	}

	return &KeyPair{
		PrivateKey:  privateKeyPEM,
		PublicKey:   publicKeyPEM,
		Fingerprint: fingerprint,
	}, nil
}
