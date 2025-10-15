.PHONY: all build build-cpp build-go clean run test install

# Variables
CXX = g++
CXXFLAGS = -std=c++17 -Wall -O2 -fPIC
INCLUDES = -I./git-core/include
LDFLAGS = -shared -lssl -lcrypto -lz
LIBDIR = ./git-core/lib
LIBNAME = libgitcore.so

# macOS specific settings
UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Darwin)
	LIBNAME = libgitcore.dylib
	LDFLAGS = -dynamiclib -lssl -lcrypto -lz
	CXXFLAGS += -I/opt/homebrew/opt/openssl/include -I/usr/local/opt/openssl/include
	LDFLAGS += -L/opt/homebrew/opt/openssl/lib -L/usr/local/opt/openssl/lib
endif

# Source files
CPP_SOURCES = $(wildcard git-core/src/*.cpp)
CPP_OBJECTS = $(CPP_SOURCES:.cpp=.o)

# Build targets
all: build

build: build-cpp build-go

build-cpp:
	@echo "Building C++ Git core library..."
	@mkdir -p $(LIBDIR)
	$(CXX) $(CXXFLAGS) $(INCLUDES) -c git-core/src/git_repository.cpp -o git-core/src/git_repository.o
	$(CXX) $(CXXFLAGS) $(INCLUDES) -c git-core/src/git_object.cpp -o git-core/src/git_object.o
	$(CXX) $(CXXFLAGS) $(INCLUDES) -c git-core/src/git_protocol.cpp -o git-core/src/git_protocol.o
	$(CXX) $(CXXFLAGS) $(INCLUDES) -c git-core/src/git_pack.cpp -o git-core/src/git_pack.o
	$(CXX) $(CXXFLAGS) $(INCLUDES) -c git-core/src/git_c_api.cpp -o git-core/src/git_c_api.o
	$(CXX) $(LDFLAGS) -o $(LIBDIR)/$(LIBNAME) \
		git-core/src/git_repository.o \
		git-core/src/git_object.o \
		git-core/src/git_protocol.o \
		git-core/src/git_pack.o \
		git-core/src/git_c_api.o
	@echo "C++ library built successfully: $(LIBDIR)/$(LIBNAME)"

build-go:
	@echo "Building Go server..."
	CGO_ENABLED=1 go build -o bin/zixiao-git-server ./cmd/server
	@echo "Go server built successfully: bin/zixiao-git-server"

clean:
	@echo "Cleaning build artifacts..."
	rm -rf git-core/src/*.o
	rm -rf $(LIBDIR)
	rm -rf bin/
	rm -rf data/*.db
	@echo "Clean complete"

run: build
	@echo "Starting ZiXiao Git Server..."
	./bin/zixiao-git-server -config ./configs/server.yaml

test:
	@echo "Running tests..."
	go test -v ./...

install:
	@echo "Installing dependencies..."
	go mod download
	@echo "Dependencies installed"

# Development helpers
dev: build
	@echo "Starting in development mode..."
	./bin/zixiao-git-server -config ./configs/server.yaml

init:
	@echo "Initializing project..."
	mkdir -p data/repositories
	mkdir -p logs
	mkdir -p bin
	mkdir -p git-core/lib
	@echo "Project initialized"

help:
	@echo "ZiXiao Git Server - Makefile"
	@echo ""
	@echo "Available targets:"
	@echo "  all        - Build everything (default)"
	@echo "  build      - Build C++ library and Go server"
	@echo "  build-cpp  - Build C++ Git core library"
	@echo "  build-go   - Build Go server"
	@echo "  clean      - Remove build artifacts"
	@echo "  run        - Build and run the server"
	@echo "  test       - Run tests"
	@echo "  install    - Install Go dependencies"
	@echo "  init       - Initialize project directories"
	@echo "  dev        - Run in development mode"
	@echo "  help       - Show this help message"
