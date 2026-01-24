DEPS = sus.go go.mod go.sum

.PHONY:
all: bin/sus bin/susd

.PHONY:
clean:
	rm -rf bin

bin/sus: sus/main.go $(DEPS) | bin
	go get
	go build -o $@ $<

bin/susd: susd/main.go $(DEPS) | bin
	go get
	go build -o $@ $<

bin:
	mkdir bin
