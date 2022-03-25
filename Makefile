RST := \033[m
BLD := \033[1m
RED := \033[31m
GRN := \033[32m
YLW := \033[33m
BLU := \033[34m
EOL := \n

all: Xcode.xip.tgz
.PHONY: all

cookies.txt:
	@printf '${BLD}${RED}make: *** [$@]${RST}${EOL}'
	@printf '${BLD}${YLW}$$${RST} '
	pipenv install
	@printf '${BLD}${YLW}$$${RST} '
	pipenv run python -m cookiebaker

ptgz:
	@printf '${BLD}${RED}make: *** [$@]${RST}${EOL}'
	@printf '${BLD}${YLW}$$${RST} '
	go get -v ./...
	@printf '${BLD}${YLW}$$${RST} '
	go build -o ./ptgz ./pseudotgz

/usr/local/bin/pseudotgz: ptgz
	@printf '${BLD}${RED}make: *** [$@]${RST}${EOL}'
	@printf '${BLD}${YLW}$$${RST} '
	cp -fv ./ptgz /usr/local/bin/pseudotgz
	@printf '${BLD}${YLW}$$${RST} '
	chmod -v +x /usr/local/bin/pseudotgz

xv:
	@printf '${BLD}${RED}make: *** [$@]${RST}${EOL}'
	@printf '${BLD}${YLW}$$${RST} '
	go get -v ./...
	@printf '${BLD}${YLW}$$${RST} '
	go build -o ./xv ./xipverifier

/usr/local/bin/xipverifier: xv
	@printf '${BLD}${RED}make: *** [$@]${RST}${EOL}'
	@printf '${BLD}${YLW}$$${RST} '
	cp -fv ./xv /usr/local/bin/xipverifier
	@printf '${BLD}${YLW}$$${RST} '
	chmod -v +x /usr/local/bin/xipverifier

Xcode.xip: cookies.txt /usr/local/bin/xipverifier
	@printf '${BLD}${RED}make: *** [$@]${RST}${EOL}'
	@printf '${BLD}${YLW}$$${RST} '
	aria2c --load-cookies='cookies.txt' --max-connection-per-server='16' --split='16' '${URL}'
	@printf '${BLD}${YLW}$$${RST} '
	xipverifier Xcode_*.xip
	@printf '${BLD}${YLW}$$${RST} '
	mv -fv Xcode_*.xip Xcode.xip

Xcode.xip.tgz: Xcode.xip /usr/local/bin/pseudotgz
	@printf '${BLD}${RED}make: *** [$@]${RST}${EOL}'
	@printf '${BLD}${YLW}$$${RST} '
	pseudotgz Xcode.xip > Xcode.xip.tgz
	@printf '${BLD}${YLW}$$${RST} '
	rm -fv Xcode.xip
