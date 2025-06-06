# Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.

include lib/make/*/Makefile

.PHONY: clencli/test
clencli/test: clencli/build go/test

.PHONY: clencli/build
clencli/build: clencli/clean go/mod/tidy go/version go/fmt go/generate go/install clencli/update-readme ## Builds the app

.PHONY: clencli/install
clencli/install: go/get go/fmt go/generate go/install ## Builds the app and install all dependencies

.PHONY: clencli/run
clencli/run: go/fmt ## Run a command
ifdef command
	make go/run command='$(command)'
else
	make go/run
endif

.PHONY: clencli/compile
clencli/compile: ## Compile to multiple architectures
	@mkdir -p dist
	@echo "Compiling for every OS and Platform"
	GOOS=darwin GOARCH=amd64 go build -o dist/clencli-darwin-amd64 main.go
	GOOS=solaris GOARCH=amd64 go build -o dist/clencli-solaris-amd64 main.go

	GOOS=freebsd GOARCH=386 go build -o dist/clencli-freebsd-386 main.go
	GOOS=freebsd GOARCH=amd64 go build -o dist/clencli-freebsd-amd64 main.go
	GOOS=freebsd GOARCH=arm go build -o dist/clencli-freebsd-arm main.go

	GOOS=openbsd GOARCH=386 go build -o dist/clencli-openbsd-386 main.go
	GOOS=openbsd GOARCH=amd64 go build -o dist/clencli-openbsd-amd64 main.go
	GOOS=openbsd GOARCH=arm go build -o dist/clencli-openbsd-arm main.go

	GOOS=linux GOARCH=386 go build -o dist/clencli-linux-386 main.go
	GOOS=linux GOARCH=amd64 go build -o dist/clencli-linux-amd64 main.go
	GOOS=linux GOARCH=arm go build -o dist/clencli-linux-arm main.go

	GOOS=windows GOARCH=386 go build -o dist/clencli-windows-386.exe main.go
	GOOS=windows GOARCH=amd64 go build -o dist/clencli-windows-amd64.exe main.go

.PHONY: clencli/clean
clencli/clean: ## Removes unnecessary files and directories
	rm -rf downloads/
	rm -rf generated-*/
	rm -rf dist/
	rm -rf build/
	rm -f box/blob.go

#rm -f $$GOPATH/bin/clencli

.PHONY: clencli/update-readme
clencli/update-readme: ## Renders template readme.tmpl with additional documents
	@echo "Updating README.tmpl to the latest version"
	@cp box/resources/init/clencli/readme.tmpl clencli/readme.tmpl
	@echo "Generate COMMANDS.md"
	@echo "## Commands" > COMMANDS.md
	@echo '```' >> COMMANDS.md
	@clencli --help >> COMMANDS.md
	@echo '```' >> COMMANDS.md
	@echo "COMMANDS.md generated successfully"
	@clencli render template --name readme

# .PHONY: clencli/test
# clencli/test: go/test

.DEFAULT_GOAL := help

.PHONY: help
help: ## This HELP message
	@fgrep -h ": ##" $(MAKEFILE_LIST) | sed -e 's/\(\:.*\#\#\)/\:\ /' | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'
