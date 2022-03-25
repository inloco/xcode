#!/bin/sh
set -ex

BIN="${1}"
LINK="$(which "${BIN}")"
REAL="$(realpath "${LINK}")"
WRAPER="/usr/local/bin/${BIN}"

cat << EOF > "${WRAPER}"
#!/bin/sh
unttify '${REAL}' "\${@}"
EOF

chmod -v +x "${WRAPER}"
