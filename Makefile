login:
	aws sso login --profile david
ecrlogin:
	aws ecr get-login-password --region eu-central-1 --profile david | docker login --username AWS --password-stdin 735968160530.dkr.ecr.eu-central-1.amazonaws.com
build:
	docker build --platform=linux/amd64 --tag 735968160530.dkr.ecr.eu-central-1.amazonaws.com/jsonconvertor:1.0 .;docker push 735968160530.dkr.ecr.eu-central-1.amazonaws.com/jsonconvertor:1.0
