build:
	go build -o dev .
run: build
	./dev init
clean:
	multipass delete dev
	multipass purge
