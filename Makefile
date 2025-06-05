

gen:
	go run github.com/99designs/gqlgen generate

dep:
	go mod tidy && go fmt

run:
	go run graphql/*.go