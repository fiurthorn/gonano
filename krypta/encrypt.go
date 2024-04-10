package krypta

import (
	"compress/gzip"
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

	armor := armor.NewWriter(f)
	defer armor.Close()

	enc, e := age.Encrypt(armor, k.recipients...)
	if e != nil {
		err = fmt.Errorf("failed to open encryption file: %w", e)
		return
	}
	defer enc.Close()

	gz := gzip.NewWriter(enc)

	if _, e := io.Copy(gz, reader); err != nil {
		err = fmt.Errorf("failed to write encrypted file: %w", e)
		return
	}
	gz.Close()
	enc.Close()

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
