
# TPP Registration Automation (Python)

This project automates the Third-Party Provider (TPP) registration flow with OP Financial Group’s (Finland) PSD2 APIs (sandbox) using Python.
It performs every step of the onboarding process — from certificate generation to TPP registration validation.

```
Under the PSD2 (Payment Services Directive 2) regulation, banks and financial institutions (known as ASPSP – Account Servicing Payment Service Provider) are required to expose secure APIs for account access and payment initiation.
To access these APIs, external entities TPPs (Third-Party Providers) must first register with the ASPSP. A TPP (Third-Party Provider) is an authorized financial service provider that connects to a bank’s API on behalf of customers.
```

# 🚀 What this project does:

Generates QWAC/QSEAL Certificates from the OP Sandbox API.

Builds and signs a Software Statement Assertion (SSA) JWT.

Generates a TPP Registration JWT, embedding the SSA.

Registers the TPP with OP’s sandbox.

Validates the registration by exchanging credentials for an access token.

Saves all generated artifacts (certs, keys, client info) for reference.

This is particularly useful for developers, QA engineers, and integrators who need to onboard TPP clients automatically in the OP PSD2 sandbox.

# TPP Registration Flow

tpp_registration/
│
├── Extracts credentials → env.json
│       Contains environment config (API keys, URLs, org info, etc.)
│
├── Calls → OP’s Certificate Generator API endpoint
│       ↓
│       Receives → QWAC & QSEAL certificates and private keys
│       ↓
│       Parses and saves them as PEM files
│
├── Builds → SSA (Software Statement Assertion)
│       A signed JWT containing TPP metadata and security info
│
├── Builds → Registration JWT
│       A signed token that embeds the SSA and is used to register the TPP
│
├── Calls → OP’s TPP Registration API endpoint
│       ↓
│       Receives → TPP client information (client_id, client_secret, etc.)
│       ↓
│       Saves → client information in JSON format
│
├── Calls → OP’s TPP Validation API endpoint
│       ↓
│       Verify → confirms TPP client registration via OAuth token exchange
│
└── ✅ Returns → TPP registration confirmation
