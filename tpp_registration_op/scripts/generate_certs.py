#!/usr/bin/env python3
import json, os, textwrap, requests
from jwcrypto import jwk # converts keys to JWK format to PEM format

def generate_certs(env):

	# Step 1: Request and save QWAC/QSEAL certificates from OP Sandbox.

	url = f"{env['sandbox_url']}?c={env['country']}&cn={env['tpp_name']}&roles={env['roles']}"
	headers = {
		"x-api-key": env["api_key"],
		"Accept": "application/json",
		"Content-Length": "0"
	}

	# Send POST request to OP Sandbox
	print("Sending request for certificates")
	response = requests.post(url, headers=headers)

	if response.status_code != 201:
		print(" Status:", response.status_code)
		raise SystemExit("Certificate generation request failed")

	data = response.json()
	print(data)

	CERT_DIR = "certs"
	os.makedirs(CERT_DIR, exist_ok=True)

	# Extract keys from response
	ssa_info = {
		"tpp_id": data.get("tppId"),
		"public_jwks_url": data.get("publicJwksUrl"),
		"qseal_kid": None,
		"alg": "RS256"
	}
	keys = data.get("privateJwks", {}).get("keys", [])

	for k in keys:
		kid = k.get("kid", "no-kid")
		k_lower = kid.lower()

		# --- Write certificate PEM from x5c[0] ---
		if "x5c" in k and k["x5c"]:
			cert_b64 = k["x5c"][0]
			cert_pem = (
				"-----BEGIN CERTIFICATE-----\n"
				+ "\n".join(textwrap.wrap(cert_b64, 64))
				+ "\n-----END CERTIFICATE-----\n"
			)

			if "qwac" in k_lower:
				open(os.path.join(CERT_DIR, "qwac_cert.pem"), "w").write(cert_pem)
			else:
				open(os.path.join(CERT_DIR, "qseal_cert.pem"), "w").write(cert_pem)
				ssa_info["qseal_kid"] = kid

		# --- Write private key PEM ---
		jwk_obj = jwk.JWK.from_json(json.dumps(k))
		try:
			priv_pem = jwk_obj.export_to_pem(private_key=True, password=None)
		except Exception as e:
			print("Error exporting private key for", kid, e)
			continue

		if "qwac" in k_lower:
			open(os.path.join(CERT_DIR, "qwac_key.pem"), "wb").write(priv_pem)
		else:
			open(os.path.join(CERT_DIR, "qseal_key.pem"), "wb").write(priv_pem)


	print("Certificates received successfully!")

	return ssa_info
