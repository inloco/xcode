RST := \033[m
BLD := \033[1m
RED := \033[31m
GRN := \033[32m
YLW := \033[33m
BLU := \033[34m
EOL := \n

all: /usr/local/bin/xipverifier cookies.txt Xcode.xip
.PHONY: all

/usr/local/bin/xipverifier: xv
	@printf '${BLD}${RED}make: *** [$@]${RST}${EOL}'
	@printf '${BLD}${YLW}$$${RST} '
	cp ./xv /usr/local/bin/xipverifier

xv:
	@printf '${BLD}${RED}make: *** [$@]${RST}${EOL}'
	@printf '${BLD}${YLW}$$${RST} '
	go get -v ./...
	@printf '${BLD}${YLW}$$${RST} '
	go build -o ./xv ./xipverifier

cookies.txt:
	@printf '${BLD}${RED}make: *** [$@]${RST}${EOL}'
	@printf '${BLD}${YLW}$$${RST} '
	pipenv install
	@printf '${BLD}${YLW}$$${RST} '
	pipenv run python -m cookiebaker

Xcode.xip:
	@printf '${BLD}${RED}make: *** [$@]${RST}${EOL}'
	@printf '${BLD}${YLW}$$${RST} '
	aria2c --load-cookies='cookies.txt' --max-connection-per-server='16' --split='16' '${URL}'
	@printf '${BLD}${YLW}$$${RST} '
	xipverifier Xcode_*.xip
	@printf '${BLD}${YLW}$$${RST} '
	mv Xcode_*.xip Xcode.xip
