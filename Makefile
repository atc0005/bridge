
# Copyright 2020 Adam Chalkley
#
# https://github.com/atc0005/bridge
#
# Licensed under the MIT License. See LICENSE file in the project root for
# full license information.

# References:
#
# https://golang.org/cmd/go/#hdr-Compile_packages_and_dependencies
# https://github.com/mapnik/sphinx-docs/blob/master/Makefile
# https://stackoverflow.com/questions/23843106/how-to-set-child-process-environment-variable-in-makefile
# https://stackoverflow.com/questions/3267145/makefile-execute-another-target
# https://unix.stackexchange.com/questions/124386/using-a-make-rule-to-call-another
# https://www.gnu.org/software/make/manual/html_node/Phony-Targets.html
# https://www.gnu.org/software/make/manual/html_node/Recipe-Syntax.html#Recipe-Syntax
# https://www.gnu.org/software/make/manual/html_node/Special-Variables.html#Special-Variables


OUTPUTDIR 				= release_assets

# https://gist.github.com/TheHippo/7e4d9ec4b7ed4c0d7a39839e6800cc16
VERSION 				= $(shell git describe --always --long --dirty)

# The default `go build` process embeds debugging information. Building
# without that debugging information reduces the binary size by around 28%.
BUILDCMD				=	go build -a -ldflags="-s -w -X main.version=${VERSION}"
GOCLEANCMD				=	go clean
GITCLEANCMD				= 	git clean -xfd
CHECKSUMCMD				=	sha256sum -b

LINTINGCMD				=   bash testing/run_linting_checks.sh
LINTINSTALLCMD			=   bash testing/install_linting_tools.sh

.DEFAULT_GOAL := help

# Targets will not work properly if a file with the same name is ever created
# in this directory. We explicitly declare our targets to be phony by
# making them a prerequisite of the special target .PHONY
.PHONY: help clean goclean gitclean pristine all windows linux linting lintinstall gotests

# WARNING: Make expects you to use tabs to introduce recipe lines
help:
	@echo "Please use \`make <target>' where <target> is one of"
	@echo "  clean          go clean to remove local build artifacts, temporary files, etc"
	@echo "  pristine       go clean and git clean local changes"
	@echo "  all            cross-compile for multiple operating systems"
	@echo "  windows        to generate a binary file for Windows"
	@echo "  linux          to generate a binary file for Linux distros"
	@echo "  lintinstall    use wrapper script to install common linting tools"
	@echo "  linting        use wrapper script to run common linting checks"
	@echo "  gotests        go test recursively, verbosely"

lintinstall:
	@echo "Calling wrapper script: $(LINTINSTALLCMD)"
	@$(LINTINSTALLCMD)
	@echo "Finished running linting tools install script"

linting:
	@echo "Calling wrapper script: $(LINTINGCMD)"
	@$(LINTINGCMD)
	@echo "Finished running linting checks"

gotests:
	@echo "Running go tests ..."
	@go test ./...
	@echo "Finished running go tests"

goclean:
	@echo "Removing object files and cached files ..."
	@$(GOCLEANCMD)
	@echo "Removing any existing release assets"
	@mkdir -p "$(OUTPUTDIR)"
	@rm -vf $(wildcard ${OUTPUTDIR}/*-linux-386)
	@rm -vf $(wildcard ${OUTPUTDIR}/*-linux-386.sha256)
	@rm -vf $(wildcard ${OUTPUTDIR}/*-linux-amd64)
	@rm -vf $(wildcard ${OUTPUTDIR}/*-linux-amd64.sha256)
	@rm -vf $(wildcard ${OUTPUTDIR}/*-windows-386.exe)
	@rm -vf $(wildcard ${OUTPUTDIR}/*-windows-386.exe.sha256)
	@rm -vf $(wildcard ${OUTPUTDIR}/*-windows-amd64.exe)
	@rm -vf $(wildcard ${OUTPUTDIR}/*-windows-amd64.exe.sha256)

# Setup alias for user reference
clean: goclean

gitclean:
	@echo "Recursively cleaning working tree by removing non-versioned files ..."
	@$(GITCLEANCMD)

pristine: goclean gitclean


# https://stackoverflow.com/questions/3267145/makefile-execute-another-target
all: clean windows linux
	@echo "Completed all cross-platform builds ..."

windows:
	@echo "Building release assets for Windows ..."

	@echo "Building 386 binary"
	@env GOOS=windows GOARCH=386 $(BUILDCMD) -o $(OUTPUTDIR)/prune-$(VERSION)-windows-386.exe ${PWD}/cmd/prune
	@env GOOS=windows GOARCH=386 $(BUILDCMD) -o $(OUTPUTDIR)/report-$(VERSION)-windows-386.exe ${PWD}/cmd/report

	@echo "Building amd64 binary"
	@env GOOS=windows GOARCH=amd64 $(BUILDCMD) -o $(OUTPUTDIR)/prune-$(VERSION)-windows-amd64.exe ${PWD}/cmd/prune
	@env GOOS=windows GOARCH=amd64 $(BUILDCMD) -o $(OUTPUTDIR)/report-$(VERSION)-windows-amd64.exe ${PWD}/cmd/report

	@echo "Generating checksum files"
	@$(CHECKSUMCMD) $(OUTPUTDIR)/prune-$(VERSION)-windows-386.exe > $(OUTPUTDIR)/prune-$(VERSION)-windows-386.exe.sha256
	@$(CHECKSUMCMD) $(OUTPUTDIR)/report-$(VERSION)-windows-386.exe > $(OUTPUTDIR)/report-$(VERSION)-windows-386.exe.sha256
	@$(CHECKSUMCMD) $(OUTPUTDIR)/prune-$(VERSION)-windows-amd64.exe > $(OUTPUTDIR)/prune-$(VERSION)-windows-amd64.exe.sha256
	@$(CHECKSUMCMD) $(OUTPUTDIR)/report-$(VERSION)-windows-amd64.exe > $(OUTPUTDIR)/report-$(VERSION)-windows-amd64.exe.sha256

	@echo "Completed build for Windows"

linux:
	@echo "Building release assets for Linux ..."
	@echo "Building 386 binary"
	@env GOOS=linux GOARCH=386 $(BUILDCMD) -o $(OUTPUTDIR)/prune-$(VERSION)-linux-386 ${PWD}/cmd/prune
	@env GOOS=linux GOARCH=386 $(BUILDCMD) -o $(OUTPUTDIR)/report-$(VERSION)-linux-386 ${PWD}/cmd/report

	@echo "Building amd64 binary"
	@env GOOS=linux GOARCH=amd64 $(BUILDCMD) -o $(OUTPUTDIR)/prune-$(VERSION)-linux-amd64 ${PWD}/cmd/prune
	@env GOOS=linux GOARCH=amd64 $(BUILDCMD) -o $(OUTPUTDIR)/report-$(VERSION)-linux-amd64 ${PWD}/cmd/report

	@echo "Generating checksum files"
	@$(CHECKSUMCMD) $(OUTPUTDIR)/prune-$(VERSION)-linux-386 > $(OUTPUTDIR)/prune-$(VERSION)-linux-386.sha256
	@$(CHECKSUMCMD) $(OUTPUTDIR)/report-$(VERSION)-linux-386 > $(OUTPUTDIR)/report-$(VERSION)-linux-386.sha256
	@$(CHECKSUMCMD) $(OUTPUTDIR)/prune-$(VERSION)-linux-amd64 > $(OUTPUTDIR)/prune-$(VERSION)-linux-amd64.sha256
	@$(CHECKSUMCMD) $(OUTPUTDIR)/report-$(VERSION)-linux-amd64 > $(OUTPUTDIR)/report-$(VERSION)-linux-amd64.sha256

	@echo "Completed build for Linux"
