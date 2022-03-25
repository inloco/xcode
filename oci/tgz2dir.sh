#!/bin/sh
set -ex

TGZPATH="$(realpath "${1}")"
TGZSIZE="$(stat -c '%s' "${TGZPATH}")"
TGZSHA256="$(shasum -a 256 "${TGZPATH}" | awk '{ print $1 }')"

TMPDIR="$(mktemp -d)"

echo 'Directory Transport Version: 1.1' > "${TMPDIR}/version"

ln -fsv "${TGZPATH}" "${TMPDIR}/${TGZSHA256}"

CONFIGPATH="${TMPDIR}/config.json"
cat << EOF > "${CONFIGPATH}"
{
  "architecture": "universal",
  "os": "darwin",
  "rootfs": {
    "type": "layers",
    "diff_ids": [
      "sha256:${TGZSHA256}"
    ]
  }
}
EOF
CONFIGSIZE="$(stat -c '%s' "${CONFIGPATH}")"
CONFIGSHA256="$(shasum -a 256 "${CONFIGPATH}" | awk '{ print $1 }')"
mv -fv "${CONFIGPATH}" "${TMPDIR}/${CONFIGSHA256}"

cat << EOF > "${TMPDIR}/manifest.json"
{
  "schemaVersion": 2,
  "config": {
    "mediaType": "application/vnd.oci.image.config.v1+json",
    "digest": "sha256:${CONFIGSHA256}",
    "size": ${CONFIGSIZE}
  },
  "layers": [
    {
      "mediaType": "application/vnd.oci.image.layer.v1.tar+gzip",
      "digest": "sha256:${TGZSHA256}",
      "size": ${TGZSIZE},
      "annotations": {
        "org.opencontainers.image.title": "${OCITIT}"
      }
    }
  ],
  "annotations": {
    "com.github.package.type": "extensible_archive",
    "org.opencontainers.image.revision": "${OCIREV}",
    "org.opencontainers.image.title": "${OCITIT}",
    "org.opencontainers.image.url": "${OCIURL}",
    "org.opencontainers.image.version": "${OCIVER}"
  }
}
EOF

rm -Rfv ./build
mv -fv "${TMPDIR}" ./build
