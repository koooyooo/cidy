.PHONY: all clean cidy install

cidy:
	@ go build -o cidy main.go

install: cidy
	@ cp cidy /usr/local/bin/cidy