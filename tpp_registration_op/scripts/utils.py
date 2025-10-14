import json
import os
import requests

def load_env(path="conf/env.json"):
	"""It loads env JSON file"""
	with open(path, "r") as env_fd:
		return json.load(env_fd)



