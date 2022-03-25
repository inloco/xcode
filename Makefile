RST := \033[m
BLD := \033[1m
RED := \033[31m
GRN := \033[32m
YLW := \033[33m
BLU := \033[34m
EOL := \n

all: Xcode.xip.tgz
	@printf '${BLD}${RED}make: *** [$@]${RST}${EOL}'
.PHONY: all

Xcode.xip: /usr/local/bin/xipverifier cookiebaker/cookies.txt
	@printf '${BLD}${RED}make: *** [$@]${RST}${EOL}'
	@printf '${BLD}${YLW}$$${RST} '
	aria2c --load-cookies='./cookiebaker/cookies.txt' --max-connection-per-server='16' --split='16' '${URL}'
	@printf '${BLD}${YLW}$$${RST} '
	xipverifier Xcode_*.xip
	@printf '${BLD}${YLW}$$${RST} '
	mv -fv Xcode_*.xip Xcode.xip

Xcode.xip.tgz: /usr/local/bin/pseudotgz Xcode.xip
	@printf '${BLD}${RED}make: *** [$@]${RST}${EOL}'
	@printf '${BLD}${YLW}$$${RST} '
	pseudotgz Xcode.xip > Xcode.xip.tgz
	@printf '${BLD}${YLW}$$${RST} '
	rm -fv Xcode.xip
