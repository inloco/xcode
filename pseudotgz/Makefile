RST := \033[m
BLD := \033[1m
RED := \033[31m
GRN := \033[32m
YLW := \033[33m
BLU := \033[34m
EOL := \n

all: build install

clean:
	@printf '${BLD}${RED}make: *** [$@]${RST}${EOL}'
	@printf '${BLD}${YLW}$$${RST} '
	rm -fv ./ptgz
.PHONY: clean

build: ptgz
	@printf '${BLD}${RED}make: *** [$@]${RST}${EOL}'
.PHONY: build

install: /usr/local/bin/pseudotgz
	@printf '${BLD}${RED}make: *** [$@]${RST}${EOL}'
.PHONY: install

ptgz:
	@printf '${BLD}${RED}make: *** [$@]${RST}${EOL}'
	@printf '${BLD}${YLW}$$${RST} '
	go get -v ./...
	@printf '${BLD}${YLW}$$${RST} '
	go build -o ./ptgz ./...

/usr/local/bin/pseudotgz: ptgz
	@printf '${BLD}${RED}make: *** [$@]${RST}${EOL}'
	@printf '${BLD}${YLW}$$${RST} '
	cp -fv ./ptgz /usr/local/bin/pseudotgz
	@printf '${BLD}${YLW}$$${RST} '
	chmod -v +x /usr/local/bin/pseudotgz
