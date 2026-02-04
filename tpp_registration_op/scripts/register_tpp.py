import requests, json


def register_tpp(env, reg_jwt):

	headers = {
		"Content-Type": "application/jwt",
		"Accept": "application/json",
		"x-api-key": env.get("api_key")
	}

	qwac_cert = ("certs/qwac_cert.pem", "certs/qwac_key.pem")

	response = requests.post(env.get("registration_url"), data=reg_jwt, headers=headers, cert=qwac_cert, verify=False)

	if response.status_code != 201:
		raise SystemExit("TPP registration request failed")

	print("✅ TPP Registration Successful. Proceeding for Validation")

	tpp_info = response.json()
	with open("certs/ttp_info.json", "w") as tpp_write:
		json.dump(tpp_info, tpp_write, indent=2)

	client_id = tpp_info.get("client_id", None)
	print(client_id)
	client_secret = tpp_info.get("client_secret", None)

	data = {
		"grant_type": "client_credentials",
		"scope": "accounts",
		"client_id": tpp_info.get("client_id", None),
		"client_secret": tpp_info.get("client_secret", None)
	}
	token_url = "https://psd2.mtls.sandbox.apis.op.fi/oauth/token"

	token_resp = requests.post(token_url, data=data, cert=qwac_cert, verify=False)
	if token_resp.status_code != 200:
		raise SystemExit("TPP validation failed")

	print("✅ TPP Validation Successful.")
