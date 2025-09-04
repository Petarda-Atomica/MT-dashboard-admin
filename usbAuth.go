package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"path"

	"github.com/shirou/gopsutil/v4/disk"
	"golang.org/x/crypto/argon2"
)

const magicHeader = "MTBLOB"

func openAndReadBlob(path string) ([]byte, error) {
	// Open file in read-only mode
	f, err := os.OpenFile(path, os.O_RDONLY, 0400)
	if err != nil {
		return nil, err
	}

	// Make sure file doesn't exceed 1 MB
	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if stat.Size() > 1<<20 {
		return nil, fmt.Errorf("blob too large")
	}

	// Copy file to buffer
	buf := make([]byte, stat.Size())
	_, err = io.ReadFull(f, buf)
	return buf, err
}

func parseBlob(blob []byte) (salt, nonce, ciphertext []byte, ok bool) {
	if len(blob) < len(magicHeader)+4+12 { // minimum: magic + saltLen + nonce
		return nil, nil, nil, false
	}

	// Check magic
	if string(blob[:len(magicHeader)]) != magicHeader {
		return nil, nil, nil, false
	}

	// Read salt length (next 4 bytes after magic)
	saltLen := int(binary.BigEndian.Uint32(blob[len(magicHeader) : len(magicHeader)+4]))

	// Ensure the blob is large enough
	if len(blob) < len(magicHeader)+4+saltLen+12 {
		return nil, nil, nil, false
	}

	// Extract salt
	saltStart := len(magicHeader) + 4
	saltEnd := saltStart + saltLen
	salt = blob[saltStart:saltEnd]

	// Extract nonce
	nonceStart := saltEnd
	nonceEnd := nonceStart + 12 // AES-GCM nonce size
	nonce = blob[nonceStart:nonceEnd]

	// Extract ciphertext
	ciphertext = blob[nonceEnd:]

	return salt, nonce, ciphertext, true
}

func decryptBlob(key, nonce, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return gcm.Open(nil, nonce, ciphertext, nil)
}

func writeBlob(passPhrase, key []byte, path string) error {
	// Generate random salt
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return err
	}

	// Derive AES key
	aesKey := argon2.IDKey(passPhrase, salt, 3, 64*1024, 2, 32)

	// Prepare to encrypt key
	block, _ := aes.NewCipher(aesKey)
	gcm, _ := cipher.NewGCM(block)
	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		return err
	}

	// Encrypt key
	ciphertext := gcm.Seal(nil, nonce, key, nil)

	// Write blob to path
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	f.Write([]byte(magicHeader))
	binary.Write(f, binary.BigEndian, uint32(len(salt)))
	f.Write(salt)
	f.Write(nonce)
	f.Write(ciphertext)

	return nil
}

func retrieveUSBKey(passPhrase []byte) (key []byte, err error) {
	blobCandidates := []string{}

	// Find all mounted partitions
	partitions, err := disk.Partitions(true)
	if err != nil {
		return nil, err
	}

	// Loop through partitions
	for _, p := range partitions {
		usage, _ := disk.Usage(p.Mountpoint)

		// Discard partitions wich are too small
		if usage == nil || usage.Total/1024/1024/1024 < 3 {
			continue
		}

		// Create candidate
		candidate := path.Join(p.Mountpoint, "wrapped.key")

		// Discard missing files
		if !fileExists(candidate) {
			continue
		}

		blobCandidates = append(blobCandidates, candidate)
	}

	// Loop through candidates
	passwordWrong := false
	for _, candidate := range blobCandidates {
		// Read blob
		blob, err := openAndReadBlob(candidate)
		if err != nil {
			log.Println("Failed to open blob:", candidate, "Reason:", err)
			continue
		}

		// Parse blob
		salt, nonce, ciphertext, ok := parseBlob(blob)
		if !ok {
			log.Println("Blob", candidate, "is not valid")
			continue
		}

		// Derive AES key from passphrase
		aesKey := argon2.IDKey(passPhrase, salt, 3, 64*1024, 2, 32)

		// Decrypt blob
		key, err := decryptBlob(aesKey, nonce, ciphertext)
		if err != nil {
			log.Println("Failed to decrypt blob", candidate, "Reason:", err)
			passwordWrong = true
			continue
		}

		return key, nil
	}

	if passwordWrong {
		return nil, fmt.Errorf("wrong password")
	}
	return nil, fmt.Errorf("no valid blobs")
}
