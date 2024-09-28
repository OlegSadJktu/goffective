# Goffective

Test project for Effective Mobile

# Run

``` sh
go mod get                                                  # get depedencies
swag init -g cmd/goffective/goffective.go --parseDependency # genearate docs
docker compose up --detach                                  # up database
go run ./cmd/goffective                                     # run project
```

