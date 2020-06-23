#!/bin/bash
set -e
#set -x

PKG_NAME=node-exporter
PKG_RPM_DIR=`pwd`/`dirname $0`
PKG_BINARY=node_exporter
PKG_VERSION=0.18.1
PKG_RELEASE=5
PKG_DISTRO=el7
PKG_CPU_ISA=x86_64
PKG_CPU_ARCH=amd64
PKG_OS=linux

PKG_RPM_FILE=${PKG_NAME}-${PKG_VERSION}-${PKG_RELEASE}.${PKG_DISTRO}.${PKG_CPU_ISA}
PKG_RPM_SPEC_FILE=${PKG_NAME}-${PKG_VERSION}-${PKG_RELEASE}.${PKG_DISTRO}.${PKG_CPU_ISA}.spec
URL_PREFIX="https://github.com/prometheus/${PKG_BINARY}/releases/download"
URL_FILENAME="${PKG_BINARY}-${PKG_VERSION}.${PKG_OS}-${PKG_CPU_ARCH}.tar.gz"
DOWNLOAD_URL="${URL_PREFIX}/v${PKG_VERSION}/${URL_FILENAME}"

printf "INFO: PKG_NAME          is set to '"$PKG_NAME"'\n";
printf "INFO: PKG_VERSION       is set to '"$PKG_VERSION"'\n";
printf "INFO: PKG_RELEASE       is set to '"$PKG_RELEASE"'\n";
printf "INFO: PKG_DISTRO        is set to '"$PKG_DISTRO"'\n";
printf "INFO: PKG_CPU_ISA       is set to '"$PKG_CPU_ISA"'\n";
printf "INFO: PKG_CPU_ARCH      is set to '"$PKG_CPU_ARCH"'\n";
printf "INFO: PKG_RPM_DIR       is set to '"${PKG_RPM_DIR}"'\n";
printf "INFO: PKG_RPM_FILE      is set to '"${PKG_RPM_FILE}"'\n";
printf "INFO: PKG_RPM_SPEC_FILE is set to '"${PKG_RPM_SPEC_FILE}"'\n";

printf "INFO: Download URL: ${DOWNLOAD_URL}\n"
printf "INFO: Download file name: ${URL_FILENAME}\n"

printf "INFO: Testing config.json file\n"
cd ${PKG_RPM_DIR}
gorpm --version
gorpm test --file config.json
if [ $? -eq 0 ]; then
    echo "INFO: Successfully validated configuration file"
else
    echo "ERROR: Failed to validate configuration file" >&2
fi
rm -rf build
mkdir -p src

if [ -f "src/${URL_FILENAME}" ]; then
  echo "INFO: src/${URL_FILENAME} exist"
else 
  echo "INFO: src/${URL_FILENAME} does not exist, downloading ..."
  curl -s -L -o src/${URL_FILENAME} ${DOWNLOAD_URL}
fi

TAR_DIR="${PKG_BINARY}-${PKG_VERSION}.${PKG_OS}-${PKG_CPU_ARCH}"
[ ! -d "src/${TAR_DIR}" ] && cd src && tar xvzf ${URL_FILENAME}


cd ${PKG_RPM_DIR}
if [ -f "src/${TAR_DIR}/${PKG_BINARY}" ]; then
  echo "INFO: src/${TAR_DIR}/${PKG_BINARY} exist"
  printf "INFO: "
  chmod +x src/${TAR_DIR}/${PKG_BINARY}
  src/${TAR_DIR}/${PKG_BINARY} --version
  mkdir -p ./usr/local/bin
  cat src/${TAR_DIR}/${PKG_BINARY} > ./usr/local/bin/${PKG_BINARY}
  chmod +x ./usr/local/bin/${PKG_BINARY}
else
  echo "ERROR: src/${TAR_DIR}/${PKG_BINARY} does not exist"
  exit 1
fi

echo "INFO: Creating directories and override macros for rpmbuild"
rm -rf ~/rpmbuild/*
mkdir -p ~/rpmbuild/{BUILD,RPMS,SOURCES,SPECS,SRPMS}
echo '%_topdir %(echo $HOME)/rpmbuild' > ~/.rpmmacros

# Architectures: https://github.com/golang/go/blob/master/src/go/build/syslist.go
# Common architectures are: amd64, 386
#
# gorpm generate-spec --version 1.0.0 --file config.json --arch 386
# gorpm generate-spec --version 1.0.0 --file config.json --arch amd64

gorpm generate-spec \
  --file config.json \
  --arch "${PKG_CPU_ARCH}" \
  --version ${PKG_VERSION} \
  --release ${PKG_RELEASE} \
  --distro ${PKG_DISTRO} \
  --cpu ${PKG_CPU_ISA} \
  --output ./spec/${PKG_RPM_SPEC_FILE}

cd ${PKG_RPM_DIR}

tar --strip-components 1 --owner=0 --group=0 \
  -czvf ~/rpmbuild/SOURCES/${PKG_RPM_FILE}.tar.gz \
  ./etc/sysconfig/${PKG_NAME}.conf \
  ./etc/profile.d/${PKG_NAME}.sh \
  ./lib/systemd/system/${PKG_NAME}.service \
  ./usr/lib/tmpfiles.d/${PKG_NAME}.conf \
  ./usr/local/bin/${PKG_BINARY}

tar -tvzf ~/rpmbuild/SOURCES/$PKG_RPM_FILE.tar.gz

rpmbuild --nodeps --target ${PKG_CPU_ISA} -ba ./spec/${PKG_RPM_SPEC_FILE}
echo "INFO: list files in ~/rpmbuild/RPMS/${PKG_CPU_ISA}/${PKG_RPM_FILE}.rpm"
rpm -qlp ~/rpmbuild/RPMS/${PKG_CPU_ISA}/${PKG_RPM_FILE}.rpm

cd ${PKG_RPM_DIR} && mkdir -p dist
rm -rf ./dist/${PKG_RPM_FILE}.rpm
cp ~/rpmbuild/RPMS/${PKG_CPU_ISA}/${PKG_RPM_FILE}.rpm ./dist/${PKG_RPM_FILE}.rpm
echo "SCP:       scp ./dist/${PKG_RPM_FILE}.rpm root@remote:/tmp/"
echo "Install:   sudo yum -y localinstall ./dist/${PKG_RPM_FILE}.rpm"
echo "RPM File:  ./dist/${PKG_RPM_FILE}.rpm"

exit 0
