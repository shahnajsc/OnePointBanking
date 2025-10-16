import requests, json


def register_tpp(env, reg_jwt):

	headers = {
		"Content-Type": "application/jwt",
		"Accept": "application/json",
		"x-api-key": env.get("api_key")
	}

	qwac_cert = ("certs/qwac_cert.pem", "certs/qwac_key.pem")

	response = requests.post(env.get("registration_url"), data=reg_jwt, headers=headers, cert=qwac_cert, verify=False)

	print("📡 Status:", response.status_code)

	resp_json = response.json()
	print(json.dumps(resp_json, indent=2))

	return response
