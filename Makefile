
.PHONY: binary clean
binary:
	mkdir -p build/ && go build -o build/main cmd/main.go

clean:
	rm -rf build testfiles/*.stdout testfiles/*.stderr
