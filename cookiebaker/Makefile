RST := \033[m
BLD := \033[1m
RED := \033[31m
GRN := \033[32m
YLW := \033[33m
BLU := \033[34m
EOL := \n

all: cookies.txt encrypt
	@printf '${BLD}${RED}make: *** [$@]${RST}${EOL}'
.PHONY: all

clean:
	@printf '${BLD}${RED}make: *** [$@]${RST}${EOL}'
	@printf '${BLD}${YLW}$$${RST} '
	rm -fv ./cookies.txt ./cookies.txt.gpg
.PHONY: clean

cookies.txt:
	@printf '${BLD}${RED}make: *** [$@]${RST}${EOL}'
	@printf '${BLD}${YLW}$$${RST} '
	cd .. && pipenv install
	@printf '${BLD}${YLW}$$${RST} '
	cd .. && pipenv run python -m cookiebaker
	@printf '${BLD}${YLW}$$${RST} '
	mv -fv ../cookies.txt .

encrypt:
	@printf '${BLD}${RED}make: *** [$@]${RST}${EOL}'
	@printf '${BLD}${YLW}$$${RST} '
	gpg --batch --no-tty --passphrase '${GPG_PASSPHRASE}' -o ./cookies.txt.gpg -c ./cookies.txt

decrypt:
	@printf '${BLD}${RED}make: *** [$@]${RST}${EOL}'
	@printf '${BLD}${YLW}$$${RST} '
	gpg --batch --no-tty --passphrase '${GPG_PASSPHRASE}' -o ./cookies.txt -d ./cookies.txt.gpg
