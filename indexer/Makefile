all:
	go mod download
	goreleaser --snapshot --rm-dist
	docker build . -t lambda-server:latest
