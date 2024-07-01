# Shell
SHELL := /bin/bash
# Out directory
OUT_DIR ?= bin
# Cmd directory
CMD_DIR ?= cmd

# name of executable (program)
execable := pdf-service

# setting colored output
NOCOLORS := 0
NOCOLORS := $(shell tput colors 2> /dev/null)
ifeq ($(shell test $(NOCOLORS) -ge 8 2> /dev/null; echo $$?), 0)
    BOLD := $(shell tput bold)
    RCOLOR := $(shell tput sgr0)
    BLACK := $(shell tput setaf 0)
    RED := $(shell tput setaf 1)
    GREEN := $(shell tput setaf 2)
    YELLOW := $(shell tput setaf 3)
    BLUE := $(shell tput setaf 4)
    MAGENTA := $(shell tput setaf 5)
    CYAN := $(shell tput setaf 6)
    WHITE := $(shell tput setaf 7)
endif

_default: _print_info _make_out_dir setup $(execable)
	@echo -e "$(GREEN)Compiled.$(RCOLOR)"

.PHONY: setup build test clean

setup:
	@echo -e "$(BOLD)$(MAGENTA)go mod download$(RCOLOR)"
	@go mod download

build:
	@$(MAKE) --no-print-directory $(MAKEFILE)

$(execable):
	@echo -e "$(BOLD)$(GREEN)go build -o $(OUT_DIR)/$(execable) $(CMD_DIR)/$(execable)/main.go $(RCOLOR)"
	@go build -o $(OUT_DIR)/$(execable) $(CMD_DIR)/$(execable)/main.go

test:
	@echo -e "$(BOLD)$(GREEN)go test $(RCOLOR)"
	go test

clean:
	@rm -rf $(OUT_DIR)
	@echo -e "$(GREEN)Cleaned.$(RCOLOR)"

_print_info:
	@echo -e "$(BLUE)Default target to make $(GREEN)$(execable)$(RCOLOR)"
	@echo -e "$(BLUE)For more info run $(MAGENTA)'make help'$(RCOLOR)"
	@echo ""

_make_out_dir:
	@mkdir -p $(OUT_DIR)

help info:
	@echo -e "\nMakefile to compile $(GREEN)$(execable)$(RCOLOR)\n"
	@echo -e "------$(CYAN) Use the following targets $(RCOLOR)-----------------"
	@echo -e "$(MAGENTA)<None>$(RCOLOR) | $(CYAN)build$(RCOLOR)\n\tto make the $(BOLD)$(GREEN)$(execable)$(RCOLOR)."
	@echo -e "$(CYAN)setup$(RCOLOR)\n\tto setup to the build."
	@echo -e "$(CYAN)test$(RCOLOR)\n\tto run tests."
	@echo -e "$(CYAN)clean$(RCOLOR)\n\tto cleanup."
	@echo -e "$(CYAN)help$(RCOLOR) | $(CYAN)info$(RCOLOR)\n\tto type this message."
	@echo -e "--------------------------------------------------\n"
	@echo -e "For more information please look into Makefile.\n"
