package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
	mathRand "math/rand"
	"time"

	"golang.org/x/crypto/sha3"
)

var letters = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func decrypt(data storedData, decryptedAAD []byte) string {
	encryptedNote := make([]byte, len(data.Note))
	copy(encryptedNote, data.Note)

	// Create new cipher and CFB instance with decrypted AAD
	blockCipher, _ := aes.NewCipher(decryptedAAD)
	cfb := cipher.NewCFBDecrypter(blockCipher, encryptedNote[:aes.BlockSize])
	// Decrypt note
	cfb.XORKeyStream(encryptedNote[aes.BlockSize:], encryptedNote[aes.BlockSize:])

	return string(encryptedNote[aes.BlockSize:])
}

func verifyNotePassword(data storedData, password string) ([]byte, error) {
	key := sha3.Sum256([]byte(password))
	blockCipher, _ := aes.NewCipher(key[:])
	decryptedAAD := make([]byte, aes.BlockSize)

	// Decrypt AAD & Hash it
	blockCipher.Decrypt(decryptedAAD, data.AADData)
	decryptedAADHash := sha3.Sum256(decryptedAAD)

	if bytes.Equal(data.AADHash[:], decryptedAADHash[:]) {
		return decryptedAAD, nil
	}
	return []byte{}, errors.New("Incorrect Password")
}

func encrypt(note string, password string) ([]byte, [32]byte, []byte) {
	// Generate additional auth data of AES blocksize
	// We cannot use 256-Bit keys for encryption
	// Without switching to CFB for AAD encryption
	AAD := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, AAD); err != nil {
		panic(err.Error())
	}

	// Create dst with length of aes blocksize + note length
	// And initialize first 16 bytes randonly for IV
	dst := make([]byte, aes.BlockSize+len([]byte(note)))
	if _, err := io.ReadFull(rand.Reader, dst[:aes.BlockSize]); err != nil {
		panic(err.Error())
	}

	// Create cipher and CFB with AAD then encrypt the note into dst
	blockCipher, _ := aes.NewCipher(AAD)
	cfb := cipher.NewCFBEncrypter(blockCipher, dst[:aes.BlockSize])
	cfb.XORKeyStream(dst[aes.BlockSize:], []byte(note))

	// Hash AAD - used for proper pass verification in get
	aadHash := sha3.Sum256(AAD)

	// Create blockCipher with hash of supplied password and encrypt AAD
	key := sha3.Sum256([]byte(password))
	blockCipher, _ = aes.NewCipher(key[:])
	blockCipher.Encrypt(AAD, AAD)

	return AAD, aadHash, dst
}

func randString(n int) string {
	mathRand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[mathRand.Intn(len(letters))]
	}
	if !storedDataEmpty(db[string(b)]) {
		return randString(n)
	}
	return string(b)
}

func storedDataEmpty(a storedData) bool {
	if bytes.Equal(a.AADData, []byte{}) {
		return true
	}
	if bytes.Equal(a.AADHash[:], []byte{}) {
		return true
	}
	if bytes.Equal(a.Note, []byte{}) {
		return true
	}
	return false
}
