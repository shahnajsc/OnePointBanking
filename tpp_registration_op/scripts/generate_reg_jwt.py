#!/usr/bin/env python3
import json, time, uuid, jwt
from datetime import datetime

def generate_reg_jwt(env, ssa_info):

	qseal_key_pem = open("certs/qseal_key.pem", "rb").read()
	time_now = int(time.time())

	headrs = {
		"alg": ssa_info.get("alg"),
		"kid": str(ssa_info.get("qseal_kid")),
		"type": "JWT"
	}

	ssa_payload = {
		"iss": ssa_info.get("tpp_id"),
		"iat": time_now,
		"exp": time_now + 3600 * 24 * 365 * 5,
		"jti": str(uuid.uuid4()),
		"software_client_id": str(uuid.uuid4()),
		"software_roles": env.get("roles", "").split(","),
		"software_jwks_endpoint": ssa_info.get("public_jwks_url"),
		"software_redirect_uris": env.get("software_redirect_uris", []),
		"software_client_name": env.get("tpp_name"),
		"software_client_uri": env.get("software_client_uri", ""),
		"org_name": env.get("tpp_name"),
		"org_id": ssa_info.get("tpp_id"),
		"org_contacts": env.get("org_contacts", [])
	}

	ssa = jwt.encode(ssa_payload, qseal_key_pem, algorithm=ssa_info.get("alg"), headers=headrs)

	open("conf/ssa.jwt", "w").write(ssa)
	print("✅ Wrote conf/ssa.jwt")

	reg_payload = {
		"iat": time_now,
		"exp": time_now + 3600,  # valid for 1 hour
		"aud": env.get("registration_url"),  # The audience (the endpoint)
		"jti": str(uuid.uuid4()),            # Unique ID for this token
		"redirect_uris": env.get("redirect_uris", []),
		"grant_types": ["client_credentials", "authorization_code", "refresh_token"],
		"software_statement": ssa  # embed the SSA you created earlier
	}

	registration_jwt = jwt.encode(reg_payload, qseal_key_pem, algorithm=ssa_info.get("alg"), headers=headrs)

	with open("conf/registration_jwt.txt", "w") as f:
		f.write(registration_jwt)
		print("✅ Wrote conf/registration_jwt.txt")

	return registration_jwt
