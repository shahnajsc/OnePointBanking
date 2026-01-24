# One Point Ledger
A PSD2/TPP compliant multi-bank financial platform.

One Point Ledger is an ongoing full-stack project designed to explore real world PSD2 (Sandbox APIs) banking integration and data driven web application development.

The goal of this project is to allow users to securely connect multiple bank accounts through a single interface, view balances, access transaction history, and eventually generate meaningful financial insights from data.

# Project Overview

The system is structured into four major components, following PSD2 Third Party Provider (TPP) architecture:
### TPP Registration & PSD2 Setup
Handles onboarding with Bank's Sandbox, registering redirect URIs, and obtaining TPP credentials under PSD2 regulations. Enables secure consent flows (AIS/PIS) and Strong Customer Authentication (SCA) required for accessing user banking data.

**Tech Stack:** Go, Sandbox APIs

### Frontend Application
A Web application where customers connect their bank accounts, view balances, and explore financial insights. Implements SCA redirect flow and future dashboards for analytics.

**Tech Stack:** React, TypeScript, TailwindCSS

### Backend API
Serves as the secure middleware layer that manages consent initiation, OAuth2 callback handling, and PSD2 token exchanges. Fetches accounts, balances, and transactions from bank's sandbox APIs and exposes clean endpoints to the frontend.

**Tech Stack:** Go, PostgreSQL, Docker

### Analytics Layer
Transforms raw PSD2 transaction history into structured financial insights using an ETL pipeline.

**Tech Stack:** (Not decided yet)

# Project Status
**TPP Registration:** Completed — onboarding with the OP Bank PSD2 Sandbox is finished (details are provided below).

**Frontend:** Actively under development.

**Backend:** Actively under development.

**Analytics Layer:** Planned.

### TPP Registration Automation (Python)

This project automates the Third-Party Provider (TPP) registration flow with OP Financial Group’s (Finland) PSD2 APIs (sandbox) using Python. It performs every step of the onboarding process — from certificate generation to TPP registration validation.

```
Under the PSD2 (Payment Services Directive 2) regulation, banks and financial institutions (known as ASPSP – Account Servicing
 Payment Service Provider) are required to expose secure APIs for account access and payment initiation. To access these APIs,
external entities TPPs (Third-Party Providers) must first register with the ASPSP. A TPP (Third-Party Provider) is an authorized
financial service provider that connects to a bank’s API on behalf of customers.
```
#### What this component does

- Generates QWAC/QSEAL Certificates from the OP Sandbox API.
- Builds and signs a Software Statement Assertion (SSA) JWT.
- Generates a TPP Registration JWT, embedding the SSA.
- Registers the TPP with OP’s sandbox.
- Validates the registration by exchanging credentials for an access token.
- Saves all generated artifacts (certs, keys, client info) for reference.

 This is particularly useful for developers, QA engineers, and integrators who need to onboard TPP clients automatically in the OP PSD2 sandbox.

#### TPP Registration Flow

```text
TPP Registration
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
│       Verify → confirms TPP client registration via authentication token exchange
│
└── Returns → TPP registration confirmation
```

#### How to Use
- Get access to OP’s PSD2 Sandbox environment.
- Create a new applicationin in OP Developer Portal to obtain your API Key.
- Create env.json file and fill your own information (check conf/env_example.json) for environment variables
- Install dependencies
  ```

  ```
- Create and activate a virtual environment
  ```

  ```
- Run the main script
  ```

  ```

