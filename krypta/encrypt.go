package krypta

import (
	"fmt"
	"io"
	"os"

	"filippo.io/age"
	"filippo.io/age/armor"
)

func (k *Krypta) EncryptFile(encrypted string, reader io.Reader) (err error) {
	f, err := os.Create(encrypted)
	if err != nil {
		err = fmt.Errorf("failed to open file to encrypt (%s): %w", encrypted, err)
		return
	}
	defer f.Close()

	armorWriter := armor.NewWriter(f)
	defer armorWriter.Close()

	w, e := age.Encrypt(armorWriter, k.recipients...)
	if e != nil {
		err = fmt.Errorf("failed to open encryption file: %w", e)
		return
	}
	defer w.Close()

	if _, e := io.Copy(w, reader); err != nil {
		err = fmt.Errorf("failed to write encrypted file: %w", e)
		return
	}

	return
}

func writeEncryptedSecret(filePath, key string, plain io.Reader) (err error) {
	of, err := os.Create(filePath)
	if err != nil {
		return
	}
	rcpt, err := age.NewScryptRecipient(key)
	if err != nil {
		return
	}
	aof := armor.NewWriter(of)
	defer aof.Close()
	w, err := age.Encrypt(aof, rcpt)
	if err != nil {
		return
	}
	defer w.Close()

	io.Copy(w, plain)
	return
}
