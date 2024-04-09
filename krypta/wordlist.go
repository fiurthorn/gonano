package krypta

import (
	"crypto/rand"
	"encoding/binary"
	"strings"

	_ "embed"
)

// Copyright 2019 Google LLC
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd
func randomWord() string {
	buf := make([]byte, 2)
	if _, err := rand.Read(buf); err != nil {
		panic(err)
	}
	n := binary.BigEndian.Uint16(buf)
	return wordlist[int(n)%2048]
}

// wordlist is the BIP39 list of 2048 english words, and it's used to generate
// the suggested passphrases.
//
//go:embed wordlist
var _wordlist string

var wordlist = strings.Split(_wordlist, "\n")
