

all:
	make template
	make css
	make build
	make run

template:
	templ generate

css: 
	tailwindcss -i ./static/css/input.css -o ./static/css/output.css

build: 
	go build

run:
	./app.exe
