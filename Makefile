
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

OUTPUTDIR 				:= release_assets

# https://gist.github.com/TheHippo/7e4d9ec4b7ed4c0d7a39839e6800cc16
VERSION 				:= $(shell git describe --always --long --dirty)

# The default `go build` process embeds debugging information. Building
# without that debugging information reduces the binary size by around 28%.
BUILDCMD				:=	go build -a -ldflags="-s -w -X main.version=${VERSION}"
GOCLEANCMD				:=	go clean
GITCLEANCMD				:= 	git clean -xfd
CHECKSUMCMD				:=	sha256sum -b

LINTINGCMD				:=   bash testing/run_linting_checks.sh
LINTINSTALLCMD			:=   bash testing/install_linting_tools.sh

.DEFAULT_GOAL := help

# Targets will not work properly if a file with the same name is ever created
# in this directory. We explicitly declare our targets to be phony by
# making them a prerequisite of the special target .PHONY
.PHONY: help clean goclean gitclean pristine all windows linux linting lintinstall gotests build-prune build-report

# WARNING: Make expects you to use tabs to introduce recipe lines
help:
	@echo "Please use \`make <target>' where <target> is one of"
	@echo "  clean          go clean to remove local build artifacts, temporary files, etc"
	@echo "  pristine       go clean and git clean local changes"
	@echo "  all            to generate binary files for Windows and Linux"
	@echo "  build-prune    to generate binary files for the prune command"
	@echo "  build-report   to generate binary files for the report command"
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
	@rm -vf $(wildcard ${OUTPUTDIR}/*/*-linux-*)
	@rm -vf $(wildcard ${OUTPUTDIR}/*/*-windows-*)

# Setup alias for user reference
clean: goclean

gitclean:
	@echo "Recursively cleaning working tree by removing non-versioned files ..."
	@$(GITCLEANCMD)

pristine: goclean gitclean

# https://stackoverflow.com/questions/3267145/makefile-execute-another-target
all: clean build-prune build-report
	@echo "Completed all cross-platform builds ..."

build-prune:
	@echo "Building release assets for prune ..."

	@mkdir -p $(OUTPUTDIR)/prune

	@echo "Building 386 binaries"
	@env GOOS=windows GOARCH=386 $(BUILDCMD) -o $(OUTPUTDIR)/prune/prune-$(VERSION)-windows-386.exe ${PWD}/cmd/prune
	@env GOOS=linux GOARCH=386 $(BUILDCMD) -o $(OUTPUTDIR)/prune/prune-$(VERSION)-linux-386 ${PWD}/cmd/prune

	@echo "Building amd64 binaries"
	@env GOOS=windows GOARCH=amd64 $(BUILDCMD) -o $(OUTPUTDIR)/prune/prune-$(VERSION)-windows-amd64.exe ${PWD}/cmd/prune
	@env GOOS=linux GOARCH=amd64 $(BUILDCMD) -o $(OUTPUTDIR)/prune/prune-$(VERSION)-linux-amd64 ${PWD}/cmd/prune

	@echo "Generating checksum files"
	@$(CHECKSUMCMD) $(OUTPUTDIR)/prune/prune-$(VERSION)-windows-386.exe > $(OUTPUTDIR)/prune/prune-$(VERSION)-windows-386.exe.sha256
	@$(CHECKSUMCMD) $(OUTPUTDIR)/prune/prune-$(VERSION)-windows-amd64.exe > $(OUTPUTDIR)/prune/prune-$(VERSION)-windows-amd64.exe.sha256
	@$(CHECKSUMCMD) "$(OUTPUTDIR)/prune/prune-$(VERSION)-linux-386" > "$(OUTPUTDIR)/prune/prune-$(VERSION)-linux-386.sha256"
	@$(CHECKSUMCMD) "$(OUTPUTDIR)/prune/prune-$(VERSION)-linux-amd64" > "$(OUTPUTDIR)/prune/prune-$(VERSION)-linux-amd64.sha256"

	@echo "Completed build for prune"

build-report:
	@echo "Building release assets for report ..."

	@mkdir -p $(OUTPUTDIR)/report

	@echo "Building 386 binaries"
	@env GOOS=windows GOARCH=386 $(BUILDCMD) -o $(OUTPUTDIR)/report/report-$(VERSION)-windows-386.exe ${PWD}/cmd/report
	@env GOOS=linux GOARCH=386 $(BUILDCMD) -o $(OUTPUTDIR)/report/report-$(VERSION)-linux-386 ${PWD}/cmd/report

	@echo "Building amd64 binaries"
	@env GOOS=windows GOARCH=amd64 $(BUILDCMD) -o $(OUTPUTDIR)/report/report-$(VERSION)-windows-amd64.exe ${PWD}/cmd/report
	@env GOOS=linux GOARCH=amd64 $(BUILDCMD) -o $(OUTPUTDIR)/report/report-$(VERSION)-linux-amd64 ${PWD}/cmd/report

	@echo "Generating checksum files"
	@$(CHECKSUMCMD) $(OUTPUTDIR)/report/report-$(VERSION)-windows-386.exe > $(OUTPUTDIR)/report/report-$(VERSION)-windows-386.exe.sha256
	@$(CHECKSUMCMD) $(OUTPUTDIR)/report/report-$(VERSION)-windows-amd64.exe > $(OUTPUTDIR)/report/report-$(VERSION)-windows-amd64.exe.sha256
	@$(CHECKSUMCMD) "$(OUTPUTDIR)/report/report-$(VERSION)-linux-386" > "$(OUTPUTDIR)/report/report-$(VERSION)-linux-386.sha256"
	@$(CHECKSUMCMD) "$(OUTPUTDIR)/report/report-$(VERSION)-linux-amd64" > "$(OUTPUTDIR)/report/report-$(VERSION)-linux-amd64.sha256"

	@echo "Completed build for report"
