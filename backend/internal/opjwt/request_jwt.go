package opjwt

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type RequestClaims struct {
	Aud          string
	Iss          string
	ClientID     string
	RedirectURI  string
	Scope        string
	State        string
	Nonce        string
	AuthorizationID string
	Kid          string
}

func LoadRSAPrivateKeyFromPEM(path string) (*rsa.PrivateKey, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read qseal key: %w", err)
	}

	block, _ := pem.Decode(b)
	if block == nil {
		return nil, fmt.Errorf("invalid PEM in %s", path)
	}

	// Try PKCS1 first
	if key, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		return key, nil
	}

	// Then PKCS8
	k, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse private key (pkcs1/pkcs8): %w", err)
	}

	rsaKey, ok := k.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("not RSA private key")
	}
	return rsaKey, nil
}

func SignOPRequestJWT(priv *rsa.PrivateKey, kid string, c RequestClaims) (string, error) {
	now := time.Now()

	// OP expects claims.authorizationId inside userinfo + id_token in their examples
	claims := jwt.MapClaims{
		"aud":           c.Aud,
		"iss":           c.Iss,
		"response_type": "code", // For backend-only; avoids fragment id_token
		"client_id":     c.ClientID,
		"redirect_uri":  c.RedirectURI,
		"scope":         c.Scope,
		"state":         c.State,
		"nonce":         c.Nonce,
		"max_age":       86400,
		"iat":           now.Unix(),
		"exp":           now.Add(5 * time.Minute).Unix(), // OP examples shows request JWT exp/iat
		"claims": map[string]any{
			"userinfo": map[string]any{
				"authorizationId": map[string]any{
					"value":     c.AuthorizationID,
					"essential": true,
				},
			},
			"id_token": map[string]any{
				"authorizationId": map[string]any{
					"value":     c.AuthorizationID,
					"essential": true,
				},
				"acr": map[string]any{
					"essential": true,
					"values": []string{
						"urn:openbanking:psd2:sca",
					},
				},
			},
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	t.Header["typ"] = "JWT"
	t.Header["kid"] = kid // qseal_kid
	return t.SignedString(priv)
}
