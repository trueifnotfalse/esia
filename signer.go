package esia

type Signer interface {
	Sign(message string) (string, error)
}
