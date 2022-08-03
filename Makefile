docker-build:
	docker image build -f Dockerfile . -t net-cat-docker
	docker images
	docker container run -p 8989:8989 -d --name net-cat-cont net-cat-docker

client-build:
	go build -o TCPChat ./client/client.go