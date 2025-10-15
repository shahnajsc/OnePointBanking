import json


def load_env(path="conf/env.json"):
	"""Load environment configuration"""
	with open(path, "r") as f:
		return json.load(f)
