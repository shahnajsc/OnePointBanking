package config

import (
	"os"
)

type Config struct {
	DatabaseURL string
	JWTSecret   string

	OPMTLSBase       string
	OPAuthBase       string
	OPClientID       string
	OPClientSecret   string
	OPAPIKey         string
	OPFAPIFinancialID string
	OPRedirectURI    string
	OPQWACCertPath   string
	OPQWACKeyPath    string
	OPQSEALKeyPath   string
	OPQSEALKid       string
}

func Load() Config {
	return Config{
		DatabaseURL:       os.Getenv("DATABASE_URL"),
		JWTSecret:         os.Getenv("JWT_SECRET"),

		OPMTLSBase:        os.Getenv("OP_MTLS_BASE"),
		OPAuthBase:        os.Getenv("OP_AUTH_BASE"),
		OPClientID:        os.Getenv("OP_CLIENT_ID"),
		OPClientSecret:    os.Getenv("OP_CLIENT_SECRET"),
		OPAPIKey:          os.Getenv("OP_API_KEY"),
		OPFAPIFinancialID: os.Getenv("OP_FAPI_FINANCIAL_ID"),
		OPRedirectURI:     os.Getenv("OP_REDIRECT_URI"),
		OPQWACCertPath:    os.Getenv("OP_QWAC_CERT_PATH"),
		OPQWACKeyPath:     os.Getenv("OP_QWAC_KEY_PATH"),
		OPQSEALKeyPath:    os.Getenv("OP_QSEAL_KEY_PATH"),
		OPQSEALKid:        os.Getenv("OP_QSEAL_KID"),
	}
}
