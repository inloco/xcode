#!/bin/sh
set -ex

function realpath {
	printf '%s\n' "$(cd -P -- "$(dirname -- "${1}")" && pwd -P)/$(basename -- "${1}")"
}

XIPPATH="$(realpath "${1}")"
XIPNAME="$(basename "${XIPPATH}")"
XIPSIZE="$(stat -f '%z' "${XIPPATH}")"
XIPSHA256="$(shasum -a 256 "${XIPPATH}" | awk '{ print $1 }')"

TMPDIR="$(mktemp -d)"

echo 'Directory Transport Version: 1.1' > "${TMPDIR}/version"

ln -s "${XIPPATH}" "${TMPDIR}/${XIPSHA256}"

CONFIGPATH="${TMPDIR}/config.json"
cat << EOF > "${CONFIGPATH}"
{
  "architecture": "universal",
  "os": "darwin",
  "rootfs": {
    "type": "layers",
    "diff_ids": [
      "sha256:${XIPSHA256}"
    ]
  }
}
EOF
CONFIGSIZE="$(stat -f '%z' "${CONFIGPATH}")"
CONFIGSHA256="$(shasum -a 256 "${CONFIGPATH}" | awk '{ print $1 }')"
mv "${CONFIGPATH}" "${TMPDIR}/${CONFIGSHA256}"

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
      "mediaType": "application/vnd.oci.image.layer.v1.xip",
      "digest": "sha256:${XIPSHA256}",
      "size": ${XIPSIZE},
      "annotations": {
        "org.opencontainers.image.title": "${XIPNAME}"
      }
    }
  ],
  "annotations": {
    "com.github.package.type": "extensible_archive",
    "org.opencontainers.image.revision": "${OCIREV}",
    "org.opencontainers.image.title": "${OCITIT}",
    "org.opencontainers.image.url": "${OCIURL}",
    "org.opencontainers.image.version": "${OCIDES}"
  }
}
EOF

rm -fR ./build
mv "${TMPDIR}" ./build
