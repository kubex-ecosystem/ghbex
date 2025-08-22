package security

// SSHKeyPair represents an SSH key pair
type SSHKeyPair struct {
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
	KeyID      int64  `json:"key_id,omitempty"`
}
