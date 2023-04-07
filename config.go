package esia

type OpenIdConfig struct {
	MnemonicsSystem    string
	RedirectUrl        string
	PortalUrl          string
	Scope              string
	CodeUrl            string
	TokenUrl           string
}

type CliSignerConfig struct {
	CertPath           string
	PrivateKeyPath     string
	PrivateKeyPassword string
	TmpPath            string
}
