package client

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"net/http"
	"sync"
	"time"

	"golang.org/x/oauth2"
)

// installationTokenSource gera e renova tokens de instalação sob demanda.
type installationTokenSource struct {
	appID          int64
	installationID int64
	privateKey     *rsa.PrivateKey
	apiBase        string
	client         *http.Client

	mu    sync.Mutex
	token *oauth2.Token
}

func (s *installationTokenSource) Token() (*oauth2.Token, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// ainda válido?
	if s.token != nil && s.token.Expiry.After(time.Now().Add(30*time.Second)) {
		return s.token, nil
	}

	// gera JWT do App
	jwt, err := makeAppJWT(s.privateKey, s.appID, 9*time.Minute)
	if err != nil {
		return nil, err
	}

	// troca por installation token
	tok, exp, err := createInstallationToken(s.client, s.apiBase, s.installationID, jwt)
	if err != nil {
		return nil, err
	}
	s.token = &oauth2.Token{
		AccessToken: tok,
		TokenType:   "token",
		Expiry:      exp,
	}
	return s.token, nil

}

func parseRSAPrivateKeyFromPEM(pemBytes []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New("no PEM block found")
	}
	switch block.Type {
	case "RSA PRIVATE KEY":
		return x509.ParsePKCS1PrivateKey(block.Bytes)
	case "PRIVATE KEY":
		k, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		rsaKey, ok := k.(*rsa.PrivateKey)
		if !ok {
			return nil, errors.New("not RSA key")
		}
		return rsaKey, nil
	default:
		return nil, errors.New("unsupported PEM type: " + block.Type)
	}
}
