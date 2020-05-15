# Setup name variables for the package/tool
NAME := bpfps
PKG := github.com/genuinetools/$(NAME)

CGO_ENABLED := 1

# Set any default go build tags.
BUILDTAGS :=

include basic.mk

.PHONY: prebuild
prebuild:
