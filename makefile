coverage:
	go test -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

coverage-linux: coverage
	xdg-open coverage.html

coverage-macos: coverage
	open coverage.html