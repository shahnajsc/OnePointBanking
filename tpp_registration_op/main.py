from scripts.utils import load_env
from scripts.generate_certs import generate_certs
from scripts.generate_reg_jwt import generate_reg_jwt
from scripts.register_tpp import register_tpp


def main():
	print("Starting OP TPP registration flow...")
	env = load_env("conf/env.json")
	ssa_data = generate_certs(env)
	reg_jwt = generate_reg_jwt(env, ssa_data)
	tpp_respons = register_tpp(env, reg_jwt)

if __name__ == "__main__":
	main()
