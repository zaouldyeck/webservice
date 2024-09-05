# Check to see if we can use ash in Alpine, or default to bash.
SHELL_PATH = /bin/ash
SHELL = $(if $(wildcard $(SHELL_PATH)),/bin/ash,/bin/bash)

run:
	go run apis/services/sales/main.go | go run apis/tooling/logfmt/main.go

# ========================================================================
# Modules support

#deps-reset:
#	git checkout -- go.mod
#	go mod tidy
#	go mod vendor

tidy:
	go mod tidy
	go mod vendor
