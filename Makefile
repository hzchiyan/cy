build:
	go build -ldflags="-s -w" -o cy main.go
	$(if $(shell command -v upx), upx cy)

mac:
	GOOS=darwin go build -ldflags="-s -w" -o cy-darwin .
	$(if $(shell command -v upx), upx cy-darwin)

win:
	GOOS=windows go build -ldflags="-s -w" -o cy.exe .
	$(if $(shell command -v upx), upx cy.exe)

linux:
	GOOS=linux go build -ldflags="-s -w" -o cy-linux .
	$(if $(shell command -v upx), upx cy-linux)