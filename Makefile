tests:
	go test -race -count=1 -coverprofile=coverage.out ./...

code-coverage:
	go tool cover -func=coverage.out

gci:
	gci write --skip-generated --skip-vendor .

lint:
	golangci-lint run

addlicense:
	addlicense -f LICENSE -ignore=vendor/\*\* -ignore=.\* -ignore=.\*/\*\*/\* .

addlicense-check:
	addlicense -f LICENSE -ignore=vendor/\*\* -ignore=.\* -ignore=.\*/\*\*/\* -check .
