# Setup name variables for the package/tool
NAME := bpfps
PKG := github.com/genuinetools/$(NAME)

CGO_ENABLED := 1

include basic.mk
