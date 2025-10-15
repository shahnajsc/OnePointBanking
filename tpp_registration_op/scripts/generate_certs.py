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
		print(response.text)
		raise SystemExit("Certificate request failed")

	print(response.headers)
	print("Certificates received successfully!")

	data = response.json()

	# Save the raw response for reference // ** delete this if not needed in future
	with open("conf/response.json", "w") as resp_file:
		json.dump(data, resp_file, indent=2)
		print("🗂️ Saved full response to response.json")

	CERT_DIR = "certs"
	os.makedirs(CERT_DIR, exist_ok=True)
	JWKS_OUT = "certs/public_jwks.json"

	# Extract keys from response
	ssa_info = {
		"tpp_id": data.get("tppId"),
		"public_jwks_url": data.get("publicJwksUrl"),
		"qseal_kid": None,
		"alg": "RS256"
	}
	keys = data.get("privateJwks", {}).get("keys", [])
	public_jwks = {"keys": []}

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
				print("🧾 Wrote certs/qwac_cert.pem")
			else:
				open(os.path.join(CERT_DIR, "qseal_cert.pem"), "w").write(cert_pem)
				print("🧾 Wrote certs/qseal_cert.pem")

		# --- Write private key PEM ---
		jwk_obj = jwk.JWK.from_json(json.dumps(k))
		try:
			priv_pem = jwk_obj.export_to_pem(private_key=True, password=None)
		except Exception as e:
			print("⚠️ Error exporting private key for", kid, e)
			continue

		if "qwac" in k_lower:
			open(os.path.join(CERT_DIR, "qwac_key.pem"), "wb").write(priv_pem)
			print("🔑 Wrote certs/qwac_key.pem")
		else:
			open(os.path.join(CERT_DIR, "qseal_key.pem"), "wb").write(priv_pem)
			print("🔑 Wrote certs/qseal_key.pem")

		# --- Add public JWK to JWKS file ---
		pub_jwk_json = jwk_obj.export(private_key=False)
		public_jwks["keys"].append(json.loads(pub_jwk_json))

		if "qseal" in k_lower:
			ssa_info["qseal_kid"] = k["kid"]

	# Save the public JWKS file
	with open(JWKS_OUT, "w") as f:
		json.dump(public_jwks, f, indent=2)
		print("🔓 Wrote", JWKS_OUT)

	print("\n✅ Certificate generation complete!")
	print("Files created in ./certs and ./conf")

	return ssa_info
