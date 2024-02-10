```bash
# install dependency
go get github.com/spf13/cobra


# run all tests
go test ./... && echo ":)"

# run all tests and show output
go test -v ./... && echo ":)"

# skip integration tests
go test -v --skip "_IT" ./... && echo ":)"
```