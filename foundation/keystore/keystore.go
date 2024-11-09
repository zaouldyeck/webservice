// Package keystore implements an in-memory keystore for JWT support.
package keystore

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"path"
	"strings"
)

// key represents key information.
type key struct {
	privatePEM string
	publicPEM  string
}

// KeyStore represents an in memory store implementation of KeyLookup interface.
type KeyStore struct {
	store map[string]key
}

// New constructs an empty KeyStore, ready to be used.
func New() *KeyStore {
	return &KeyStore{
		store: make(map[string]key),
	}
}

// LoadRSAKeys loads a set of RSA pems. Name of pem file will be used as key id.
func (ks *KeyStore) LoadRSAKeys(fsys fs.FS) error {
	fn := func(filename string, dirEntry fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("walkdir failed: %w", err)
		}

		if dirEntry.IsDir() {
			return nil
		}

		if path.Ext(filename) != ".pem" {
			return nil
		}

		file, err := fsys.Open(filename)
		if err != nil {
			return fmt.Errorf("open key file failed: %w", err)
		}
		defer file.Close()

		// Limit pem file size to 1 MB.
		pem, err := io.ReadAll(io.LimitReader(file, 1024*1024))
		if err != nil {
			return fmt.Errorf("read auth private key file failed: %w", err)
		}

		privatePEM := string(pem)
		publicPEM, err := toPublicPEM(privatePEM)
		if err != nil {
			return fmt.Errorf("convert private key file to public failed: %w", err)
		}

		key := key{
			privatePEM: privatePEM,
			publicPEM:  publicPEM,
		}

		ks.store[strings.TrimSuffix(dirEntry.Name(), ".pem")] = key

		return nil
	}

	if err := fs.WalkDir(fsys, ".", fn); err != nil {
		return fmt.Errorf("walkdir failed: %w", err)
	}

	return nil
}

func toPublicPEM(privatePEM string) (string, error) {
	block, _ := pem.Decode([]byte(privatePEM))
	if block == nil {
		return "", errors.New("invalid key: Key must be a PEM encoded PKCS1 or PKCS8 key")
	}

	var parsedKey any
	parsedKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		parsedKey, err = x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return "", err
		}
	}

	pk, ok := parsedKey.(*rsa.PrivateKey)
	if !ok {
		return "", errors.New("key is not a valid RSA private key")
	}

	asn1Bytes, err := x509.MarshalPKIXPublicKey(&pk.PublicKey)
	if err != nil {
		return "", fmt.Errorf("marshaling public key: %w", err)
	}

	publicBlock := pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: asn1Bytes,
	}

	var buf bytes.Buffer
	if err := pem.Encode(&buf, &publicBlock); err != nil {
		return "", fmt.Errorf("encoding to public PEM: %w", err)
	}

	return buf.String(), nil
}
