package krypta

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"filippo.io/age"
	"filippo.io/age/armor"
)

func (k *Krypta) DecryptFile(encrypted string) (content string, err error) {
	if _, e := os.Stat(encrypted); os.IsNotExist(e) {
		return
	} else if e != nil {
		err = e
		return
	}

	f, err := os.Open(encrypted)
	if err != nil {
		err = fmt.Errorf("failed to open file to decrypt (%s): %w", encrypted, err)
		return
	}
	defer f.Close()

	armorReader := armor.NewReader(f)
	r, err := age.Decrypt(armorReader, k.identities...)
	if err != nil {
		err = fmt.Errorf("failed to open encrypted file: %v", err)
		return
	}

	out := &bytes.Buffer{}
	if _, e := io.Copy(out, r); e != nil {
		err = fmt.Errorf("failed to read encrypted file: %w", e)
		return
	}

	content = out.String()
	return
}
