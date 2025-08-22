package ghclient

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func makeAppJWT(priv *rsa.PrivateKey, appID int64, dur time.Duration) (string, error) {
	now := time.Now()
	iat := now.Add(-30 * time.Second).Unix()
	exp := now.Add(dur).Unix()

	header := map[string]any{"alg": "RS256", "typ": "JWT"}
	payload := map[string]any{"iat": iat, "exp": exp, "iss": fmt.Sprintf("%d", appID)}

	hb, _ := json.Marshal(header)
	pb, _ := json.Marshal(payload)

	enc := func(b []byte) string { return base64.RawURLEncoding.EncodeToString(b) }

	unsigned := enc(hb) + "." + enc(pb)
	sig, err := rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA256, sha256Bytes(unsigned))
	if err != nil {
		return "", err
	}
	return unsigned + "." + enc(sig), nil
}

func sha256Bytes(s string) []byte {
	h := crypto.SHA256.New()
	h.Write([]byte(s))
	return h.Sum(nil)
}

func createInstallationToken(hc *http.Client, base string, instID int64, jwt string) (token string, expiry time.Time, err error) {
	url := fmt.Sprintf("%s/app/installations/%d/access_tokens", base, instID)
	req, _ := http.NewRequest("POST", url, bytes.NewReader([]byte("{}")))
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "Bearer "+jwt)
	resp, err := hc.Do(req)
	if err != nil {
		return "", time.Time{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		return "", time.Time{}, fmt.Errorf("installation token http %d", resp.StatusCode)
	}
	var out struct {
		Token     string    `json:"token"`
		ExpiresAt time.Time `json:"expires_at"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", time.Time{}, err
	}
	return out.Token, out.ExpiresAt, nil
}
