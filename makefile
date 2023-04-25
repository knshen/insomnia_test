mock:
	rm -rf mocks/*
	mockgen -source=client/auth.go -package=mockclient -destination=mocks/mockclient/auth.go
	mockgen -source=client/project.go -package=mockclient -destination=mocks/mockclient/project.go

gomod:
	go mod tidy

run:
	go build code.sk.org/insomnia_test
	./insomnia_test

all:
	make gomod
	make run
