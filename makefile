

all:
	make template
	make build
	make run

template:
	templ generate

build: 
	go build -o app.exe

run:
	./app.exe
