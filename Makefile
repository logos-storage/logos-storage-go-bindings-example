# Destination folder for the downloaded libraries
LIBS_DIR := $(abspath ./libs)

# Flags for CGO to find the headers and the shared library
UNAME_S := $(shell uname -s)
CGO_CFLAGS  := -I$(LIBS_DIR)
CGO_LDFLAGS := -L$(LIBS_DIR) -lstorage -Wl,-rpath,$(LIBS_DIR)

ifeq ($(OS),Windows_NT)
  BIN_NAME := example.exe
else
  BIN_NAME := example
endif

# Configuration for fetching the right binary
OS ?= "linux"
ARCH ?= "amd64"
VERSION ?= "v0.3.1"
DOWNLOAD_URL := "https://github.com/logos-storage/logos-storage-nim/releases/download/$(VERSION)/libstorage-${OS}-${ARCH}-$(VERSION).zip"		

all: run

fetch: 
	@echo "Fetching libstorage from GitHub Actions from: ${DOWNLOAD_URL}"
	curl -fSL --create-dirs -o $(LIBS_DIR)/libstorage-${OS}-${ARCH}-$(VERSION).zip ${DOWNLOAD_URL}
	unzip -o -qq $(LIBS_DIR)/libstorage-${OS}-${ARCH}-$(VERSION).zip -d $(LIBS_DIR)
	rm -f $(LIBS_DIR)/*.zip

build:
	CGO_ENABLED=1 CGO_CFLAGS="$(CGO_CFLAGS)" CGO_LDFLAGS="$(CGO_LDFLAGS)" go build -o $(BIN_NAME) main.go

run:
ifeq ($(OS),Windows_NT)
	pwsh -File $(CURDIR)/.github/scripts/run-windows.ps1 -BinaryName $(BIN_NAME)
else
	./$(BIN_NAME)
endif

clean:
	rm -f $(BIN_NAME)
	rm -Rf $(LIBS_DIR)/*