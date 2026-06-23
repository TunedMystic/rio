APP="$$(basename -- $$(PWD))"



# -------------------------------------------------------------------
# App-related commands
# -------------------------------------------------------------------

## @(app) - 🧪 Run tests
test:
	@echo "✨🧪✨ Running tests\n"
	@go test $$(go list ./... | grep -v examples) -count=1


## @(app) - 🧪 Run tests and display coverage
test-v: clean
	@echo "✨🧪✨ Running tests\n"
	@go test $$(go list ./... | grep -v examples) -count=1 -covermode=atomic -coverprofile coverage.out
	@go tool cover -func coverage.out


## @(app) - 🏎️ Run benchmarks
bench:
	@echo "✨🧪✨ Running benchmarks\n"
	@go test -bench=. -benchmem -run=^$$ -count=1 ./...


## @(app) - 🎨 Serve the rio/ui component gallery on :8080
preview:
	@echo "✨📦✨ Serving component gallery on :8080\n"
	@go run ./cmd/preview


## @(app) - ✨ Remove temp files and dirs
clean:
	@echo "✨✨ Cleaning temp files\n"
	@rm -f coverage.out
	@go clean -testcache
	@find . -name '.DS_Store' -type f -delete
	@bash -c "mkdir -p bin && cd bin && find . ! -name 'watchexec' ! -name 'cwebp' ! -name 'tailwind' -type f -exec rm -f {} +"



# -------------------------------------------------------------------
# Self-documenting Makefile targets - https://bit.ly/32lE64t
# -------------------------------------------------------------------

.DEFAULT_GOAL := help

help:
	@echo "Usage:"
	@echo "  make <target>"
	@echo ""
	@echo "Targets:"
	@awk '/^[a-zA-Z\-\_0-9]+:/ \
		{ \
			helpMessage = match(lastLine, /^## (.*)/); \
			if (helpMessage) { \
				helpCommand = substr($$1, 0, index($$1, ":")-1); \
				helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
				helpGroup = match(helpMessage, /^@([^ ]*)/); \
				if (helpGroup) { \
					helpGroup = substr(helpMessage, RSTART + 1, index(helpMessage, " ")-2); \
					helpMessage = substr(helpMessage, index(helpMessage, " ")+1); \
				} \
				printf "%s|  %-20s %s\n", \
					helpGroup, helpCommand, helpMessage; \
			} \
		} \
		{ lastLine = $$0 }' \
		$(MAKEFILE_LIST) \
	| sort -t'|' -sk1,1 \
	| awk -F '|' ' \
			{ \
			cat = $$1; \
			if (cat != lastCat || lastCat == "") { \
				if ( cat == "0" ) { \
					print "\nTargets:" \
				} else { \
					gsub("_", " ", cat); \
					printf "\n%s\n", cat; \
				} \
			} \
			print $$2 \
		} \
		{ lastCat = $$1 }'
	@echo ""
