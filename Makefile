# Go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# Source paths
SOURCE_SERVER=./cmd/server
SOURCE_CLIENT=./cmd/client

# Binary path and names
BINARY_PATH=bin
BINARY_NAME_SERVER=redirect
BINARY_NAME_CLIENT=client

# Build flags
BUILD_VERSION=`git rev-parse --short HEAD`
LDFLAGS=-ldflags "-X main.Build=$(BUILD_VERSION)"

# Docker 
DOCKER_TAG=flo80/redirect

run_server: build_server
	$(BINARY_PATH)/$(BINARY_NAME_SERVER) -config testdata/redirects.json -admin localhost

all: test build 

test: 
	$(GOTEST) -v ./...

build: build_server build_client 
build_server:
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_PATH)/$(BINARY_NAME_SERVER) -v $(SOURCE_SERVER)
build_client:
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_PATH)/$(BINARY_NAME_CLIENT) -v $(SOURCE_CLIENT) 
 
clean:
	$(GOCLEAN)
	rm -f $(BINARY_PATH)/$(BINARY_NAME_SERVER)
	rm -f $(BINARY_PATH)/$(BINARY_NAME_CLIENT)


# Cross compilation
build_all_architectures: build_linux build_mac build_raspbian3 build_windows

build_linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_PATH)/linux/$(BINARY_NAME_SERVER) -v $(SOURCE_SERVER)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_PATH)/linux/$(BINARY_NAME_CLIENT) -v $(SOURCE_CLIENT) 

build_raspbian3:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=5 $(GOBUILD) $(LDFLAGS) -o $(BINARY_PATH)/raspbian3/$(BINARY_NAME_SERVER) -v $(SOURCE_SERVER)
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=5 $(GOBUILD) $(LDFLAGS) -o $(BINARY_PATH)/raspbian3/$(BINARY_NAME_CLIENT) -v $(SOURCE_CLIENT) 

build_windows:
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_PATH)/windows/$(BINARY_NAME_SERVER).exe -v $(SOURCE_SERVER)
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_PATH)/windows/$(BINARY_NAME_CLIENT).exe -v $(SOURCE_CLIENT) 

build_mac:
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_PATH)/mac/$(BINARY_NAME_SERVER) -v $(SOURCE_SERVER)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_PATH)/mac/$(BINARY_NAME_CLIENT) -v $(SOURCE_CLIENT) 

clean_cross_compilation:
	rm -f $(BINARY_PATH)/linux/$(BINARY_NAME_SERVER)
	rm -f $(BINARY_PATH)/linux/$(BINARY_NAME_CLIENT)
	rm -f $(BINARY_PATH)/raspbian3/$(BINARY_NAME_SERVER)
	rm -f $(BINARY_PATH)/raspbian3/$(BINARY_NAME_CLIENT)
	rm -f $(BINARY_PATH)/windows/$(BINARY_NAME_SERVER).exe
	rm -f $(BINARY_PATH)/windows/$(BINARY_NAME_CLIENT).exe
	rm -f $(BINARY_PATH)/mac/$(BINARY_NAME_SERVER)
	rm -f $(BINARY_PATH)/mac/$(BINARY_NAME_CLIENT)


docker_all: docker_clean docker docker_test

docker:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD)  -ldflags "-X main.Build=$(BUILD_VERSION) -s" -o $(BINARY_PATH)/docker/$(BINARY_NAME_SERVER) -v $(SOURCE_SERVER) 
	docker build -t $(DOCKER_TAG) .
	
docker_test:
	docker run -v `pwd`/testdata:/redirects -it -p 80:80 --rm $(DOCKER_TAG)

docker_clean:
	rm -f $(BINARY_PATH)/docker/$(BINARY_NAME_SERVER)
