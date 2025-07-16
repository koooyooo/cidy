.PHONY: all clean cidy install

cidy:
	@ go build -o cidy main.go

install: cidy
	@ sudo cp cidy /usr/local/bin/cidy