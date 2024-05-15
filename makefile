build: build-windows build-linux build-darwin

build-windows:
	GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o ./build/windows/githubrunner.exe

build-linux:
	GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o ./build/linux/githubrunner

build-darwin:
	GOOS=darwin GOARCH=arm64 go build -ldflags "-s -w" -o ./build/darwin/githubrunner

build-docker:
	sudo docker build -t "miras-github-runner:alpha" .