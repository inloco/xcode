#!/bin/sh
set -ex

CRTDIR=/usr/local/share/ca-certificates/apple
mkdir -p "${CRTDIR}"

BASE_URL=https://www.apple.com
for CRT in AppleIncRootCertificate AppleComputerRootCertificate AppleRootCA-G2 AppleRootCA-G3
do
  wget -P "${CRTDIR}" "${BASE_URL}/appleca/${CRT}.cer" || wget -P "${CRTDIR}" "${BASE_URL}/certificateauthority/${CRT}.cer"
  openssl x509 -inform DER -in "${CRTDIR}/${CRT}.cer" -out "${CRTDIR}/${CRT}.crt"
  rm -f "${CRTDIR}/${CRT}.cer"
done

update-ca-certificates
