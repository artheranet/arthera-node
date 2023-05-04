#!/usr/bin/env bash

SCRIPT_DIR=$(cd "$(dirname "$0")" && pwd -P)

if [[ -z "${VERSION}" ]]; then
  echo "VERSION is not defined"
  exit 1
fi

S3_DEST="s3://release.arthera.net"

cp -f "$SCRIPT_DIR/installer" /tmp/installer
sed -i "s/ARTHERA_RELEASE=0.0.0/ARTHERA_RELEASE=$VERSION/g" /tmp/installer
aws s3 cp /tmp/installer "$S3_DEST/$VERSION/installer"
rm -f /tmp/installer

DEST_LOCATION="$S3_DEST/$VERSION/arthera"
echo "Uploading 'arthera-node' to $DEST_LOCATION"
aws s3 cp "$SCRIPT_DIR/../build/arthera-node" "$DEST_LOCATION"

ETH_BINARIES=('bootnode' 'abidump' 'abigen')
for binary in "${ETH_BINARIES[@]}"
do
	DEST_LOCATION="$S3_DEST/$VERSION/$binary"
	echo "Uploading '$binary' to $DEST_LOCATION"
	aws s3 cp "$SCRIPT_DIR/../../arthera-go-ethereum/build/bin/$binary" $DEST_LOCATION
done

aws cloudfront create-invalidation --distribution-id E2UB9BK34N02BW --paths "/*"
exit 0
