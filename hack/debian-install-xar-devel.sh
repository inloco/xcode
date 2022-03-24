#!/bin/sh
set -ex

TMPDIR=$(mktemp -d)

mkdir -p "${TMPDIR}/usr"
ln -s /usr/include "${TMPDIR}/usr/include"
ln -s /usr/lib/x86_64-linux-gnu "${TMPDIR}/usr/lib64"

for PKG in xar xar-devel
do
  wget -O "${TMPDIR}/${PKG}.rpm" "https://download-ib01.fedoraproject.org/pub/epel/8/Everything/x86_64/Packages/x/${PKG}-1.8.0.417.1-2.el8.x86_64.rpm"
  rpm2cpio "${TMPDIR}/${PKG}.rpm" > "${TMPDIR}/${PKG}.cpio"
  cpio -vmidD "${TMPDIR}" < "${TMPDIR}/${PKG}.cpio"
done

rm -fR "${TMPDIR}"
