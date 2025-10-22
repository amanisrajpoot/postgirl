APP=litepost
VERSION?=$(shell git describe --tags --always --dirty)
LDFLAGS=-s -w -X main.version=$(VERSION)
BUILD_CGO=CGO_ENABLED=1 go build -trimpath -ldflags "$(LDFLAGS)"
BUILD_STATIC=CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)"

DIST=dist

PLATFORMS=\
  darwin/amd64 \
  darwin/arm64 \
  linux/amd64  \
  linux/arm64  \
  windows/amd64 \
  windows/arm64

all: clean build-all

clean:
	rm -rf $(DIST)

build-%:
	@os=$(word 1,$(subst /, ,$*)) ; arch=$(word 2,$(subst /, ,$*)); \
	ext=""; [ "$$os" = "windows" ] && ext=".exe"; \
	out="$(DIST)/$(APP)-$$os-$$arch$$ext"; \
	mkdir -p $$(dirname "$$out"); \
	if [ "$$os" = "darwin" ]; then \
		GOOS=$$os GOARCH=$$arch $(BUILD_CGO) -o "$$out" ./cmd/$(APP); \
	else \
		GOOS=$$os GOARCH=$$arch $(BUILD_STATIC) -o "$$out" ./cmd/$(APP); \
	fi

build-all: $(addprefix build-,$(PLATFORMS))

.PHONY: all clean build-all $(addprefix build-,$(PLATFORMS))
