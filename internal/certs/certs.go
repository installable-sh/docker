package certs

import _ "embed"

//go:embed ca-certificates.crt
var CACerts []byte
