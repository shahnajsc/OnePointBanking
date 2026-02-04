package service

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"net/url"
	"time"

	"github.com/shahnajsc/OnePointLedger/backend/internal/opclient"
	"github.com/shahnajsc/OnePointLedger/backend/internal/opjwt"
	"github.com/shahnajsc/OnePointLedger/backend/internal/repo"
)

type OPConnectService struct {
	op          *opclient.AISClient
	repo        *repo.OPConnectRepo
	authBase    string
	redirectURI string
	clientID    string
	aud         string
	qsealKey    *rsa.PrivateKey
	qsealKid    string
}

func NewOPConnectService(
	op *opclient.AISClient,
	repo *repo.OPConnectRepo,
	authBase, redirectURI, clientID, aud, qsealKeyPath, qsealKid string,
) (*OPConnectService, error) {
	priv, err := opjwt.LoadRSAPrivateKeyFromPEM(qsealKeyPath)
	if err != nil {
		return nil, err
	}
	return &OPConnectService{
		op:          op,
		repo:        repo,
		authBase:    authBase,
		redirectURI: redirectURI,
		clientID:    clientID,
		aud:         aud,
		qsealKey:    priv,
		qsealKid:    qsealKid,
	}, nil
}

func (s *OPConnectService) Start(ctx context.Context, userID string) (string, error) {
	// Client credentials token
	tctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	ccToken, err := s.op.ClientCredentialsToken(tctx)
	if err != nil {
		return "", fmt.Errorf("client credentials token: %w", err)
	}

	// Create authorization intent
	actx, cancel2 := context.WithTimeout(ctx, 10*time.Second)
	defer cancel2()

	authorizationID, err := s.op.CreateAuthorization(actx, ccToken)
	if err != nil {
		return "", fmt.Errorf("create authorization: %w", err)
	}

	// Generate state + nonce
	state, err := randomURLSafe(24)
	if err != nil {
		return "", fmt.Errorf("state: %w", err)
	}
	nonce, err := randomURLSafe(24)
	if err != nil {
		return "", fmt.Errorf("nonce: %w", err)
	}

	// Sign OP request JWT (using TPP QSEAL RS256 algo)
	requestJWT, err := opjwt.SignOPRequestJWT(s.qsealKey, s.qsealKid, opjwt.RequestClaims{
		Aud:            s.aud,
		Iss:            s.clientID,
		ClientID:       s.clientID,
		RedirectURI:    s.redirectURI,
		Scope:          "openid accounts",
		State:          state,
		Nonce:          nonce,
		AuthorizationID: authorizationID,
	})
	if err != nil {
		return "", fmt.Errorf("sign request jwt: %w", err)
	}

	// Store pending state in Database
	sctx, cancel3 := context.WithTimeout(ctx, 3*time.Second)
	defer cancel3()

	if err := s.repo.SavePending(sctx, state, userID, authorizationID, nonce); err != nil {
		return "", fmt.Errorf("save pending authorization: %w", err)
	}

	// Build redirect URL
	u, _ := url.Parse(s.authBase + "/oauth/authorize")
	q := u.Query()
	q.Set("request", requestJWT)
	q.Set("response_type", "code")
	q.Set("client_id", s.clientID)
	q.Set("scope", "openid accounts")
	u.RawQuery = q.Encode()

	return u.String(), nil
}

func randomURLSafe(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
