package krypta

import "filippo.io/age"

type Krypta struct {
	identities []age.Identity
	recipients []age.Recipient
}

func New() *Krypta {
	return &Krypta{
		identities: ReadIdentityFile(IdentitiesFile()),
		recipients: ReadRecipient(RecipientsFile()),
	}
}
