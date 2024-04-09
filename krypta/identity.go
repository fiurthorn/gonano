package krypta

import (
	"io"
	"log"
	"os"

	"filippo.io/age"
	"filippo.io/age/armor"
)

func ReadIdentityFile(filePath string) []age.Identity {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Failed to open private keys file: %v", err)
	}

	passphrase, err := readSecret("Enter passphrase: ")
	if err != nil {
		log.Fatalf("Passphrase: %v", err)
	}

	plain := readEncryptedSecret(string(passphrase), file)

	identities, err := age.ParseIdentities(plain)
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}
	return identities
}

func ReadRecipient(filePath string) []age.Recipient {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Failed to open private keys file: %v", err)
	}
	recipient, err := age.ParseRecipients(file)
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}
	return recipient
}

func readEncryptedSecret(key string, ainf io.Reader) io.Reader {
	aof := armor.NewReader(ainf)
	idty, err := age.NewScryptIdentity(key)
	if err != nil {
		panic(err)
	}
	w, err := age.Decrypt(aof, idty)
	if err != nil {
		panic(err)
	}
	return w
}
