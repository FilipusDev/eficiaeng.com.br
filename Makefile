.PHONY: setup templ server tailwind db-dev-down db-dev-up dev-up \
	local-down local-up \
	help

.DEFAULT_GOAL := help

# --- Makefile VARS ---

export TEMPL_VERSION=v0.3.898
export TAILWINDCSS_VERSION=v4.1.10
export DAISYUI_VERSION=v5.0.43

# -- Environment detection: TOP --

ifeq (,$(filter $(MAKECMDGOALS),help setup))
	ifeq (,$(filter $(MAKECMDGOALS),db-dev-down local-down))
		ifneq (,$(filter $(MAKECMDGOALS),db-dev-up dev-up))
			# Development environment
			include .dev.env
			export $(shell sed 's/=.*//' .dev.env 2>/dev/null)
			export GIT_HASH := $(shell git rev-parse --short HEAD)
		else ifneq (,$(filter $(MAKECMDGOALS),local-up))
			# Local environment
			include .env
			export $(shell sed 's/=.*//' .env 2>/dev/null)
			export GIT_HASH := $(shell git rev-parse --short HEAD)
		endif
	endif
endif

# -- Environment detection: BOT --

# -- Dev Env: TOP --

setup: ## Check tools' installs and assets, installs and download if needed
	@if ! command -v go >/dev/null 2>&1; then \
	echo "go not found, please install..."; \
		else \
		echo "go is already installed: $$(go version)"; \
		fi
	@if ! command -v templ >/dev/null 2>&1; then \
		echo "templ not found, installing..."; \
		go install github.com/a-h/templ/cmd/templ@${TEMPL_VERSION}; \
		else \
		echo "templ is already installed: $$(templ version)"; \
		fi
	@if ! command -v air >/dev/null 2>&1; then \
		echo "air not found, installing..."; \
		go install github.com/air-verse/air@latest; \
		else \
		echo "air is already installed: $$(air -v)"; \
		fi
	@if ! command tailwindcss >/dev/null 2>&1; then \
		echo "tailwindcss not found, installing..."; \
		curl -sLo tailwindcss https://github.com/tailwindlabs/tailwindcss/releases/download/${TAILWINDCSS_VERSION}/tailwindcss-linux-x64; \
		chmod +x tailwindcss; \
		sudo mv tailwindcss /usr/local/bin/; \
		else \
		echo "tailwindcss is already installed: $$(tailwindcss -v)"; \
		fi
	@mkdir -p ./assets/css && \
	for file in daisyui.js daisyui-theme.js; do \
		if [ ! -f "./assets/css/$$file" ]; then \
			echo "$$file does not exist, downloading..."; \
			curl -sL "https://github.com/saadeghi/daisyui/releases/download/${DAISYUI_VERSION}/$$file" -o "./assets/css/$$file"; \
		else \
			echo "$$file exists"; \
		fi; \
	done

templ: ## Run templ generation in watch mode
	templ generate \
		--watch \
		--proxy="http://localhost:${APP_PORT}" \
		--open-browser=false

server: ## Run air for Go hot reload
	air \
	--build.cmd "go build -o tmp/bin/main ./cmd/web/main.go" \
	--build.full_bin "DATABASE_URL=${DATABASE_URL} tmp/bin/main" \
	--build.delay "100" \
	--build.include_ext "go" \
	--build.stop_on_error "false" \
	--misc.clean_on_exit true

tailwind: ## Watch Tailwind CSS changes
	tailwindcss -i ./assets/css/input.css -o ./assets/css/output.css --watch

clean: ## Cleans generated files from Templ and TailwindCss
	rm -f ./assets/css/output.css
	find ./templates -type f -name '*_templ.go' -delete
	sleep 3

db-dev-down: ## Stops db-dev (from db-dev-docker-compose.yml)
	docker compose --file db-dev-docker-compose.yml down --volumes --rmi local
	sleep 3

db-dev-up: ## Starts db-dev (from db-dev-docker-compose.yml)
	docker compose --file db-dev-docker-compose.yml up --build --detach
	sleep 3

dev-down: clean ## Stop development server with all watchers

dev-up: dev-down ## Start development server with all watchers
	make -j3 tailwind templ server

# -- Dev Env: BOT --

# -- Local Build: TOP --

# TODO: missing Cloudflare tunnel crap

local-down: ## Stops local server
	docker compose down --volumes --rmi local
	sleep 3

local-up: local-down ## Starts local server
	docker compose up --build --detach
	sleep 3

# -- Local Build: BOT --

# TODO: -- Prod Build: TOP --
# -- Prod Build: BOT --

# --- Help ---

help: ## Shows this help message
	@echo "Available Targets:"
	@awk 'BEGIN {FS = ":.*?## "}; \
	     /^[a-zA-Z0-9_-]+:.*?## / { \
	       if ($$1 != "Makefile") { \
	         printf "\033[36m%-20s\033[0m %s\n", $$1, $$2 \
	       } \
	     }' $(MAKEFILE_LIST) | sort
