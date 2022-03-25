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
	rm -fv ./xv
.PHONY: clean

build: xv
	@printf '${BLD}${RED}make: *** [$@]${RST}${EOL}'
.PHONY: build

install: /usr/local/bin/xipverifier
	@printf '${BLD}${RED}make: *** [$@]${RST}${EOL}'
.PHONY: install

xv:
	@printf '${BLD}${RED}make: *** [$@]${RST}${EOL}'
	@printf '${BLD}${YLW}$$${RST} '
	go get -v ./...
	@printf '${BLD}${YLW}$$${RST} '
	go build -o ./xv ./...

/usr/local/bin/xipverifier: xv
	@printf '${BLD}${RED}make: *** [$@]${RST}${EOL}'
	@printf '${BLD}${YLW}$$${RST} '
	cp -fv ./xv /usr/local/bin/xipverifier
	@printf '${BLD}${YLW}$$${RST} '
	chmod -v +x /usr/local/bin/xipverifier
