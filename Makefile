default: build

build:
	@go build .

install: build
	@mv golinkwrite $(GOPATH)/bin/golinkwrite

clean:
	@go clean
	@rm golinkwrite

all: build install
