#!/bin/sh
set -e

ARCH="x86_64"
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
EXT=".tar.gz"

ORG="puppetlabs"
REPO="pdkgo"
APP="pct"
APP_PKG_NAME="pct"

NO_TEL=${1:-false}

if [ ${NO_TEL} = "--no-telemetry" ]; then
  APP_PKG_NAME="notel_pct"
fi

RELEASES=""
FILE=""
CHECKSUM=""

logDebug() {
  if [ ! -z $PCT_INSTALL_DEBUG ]; then
    echo $1
  fi
}

getChecksums() {
  for i in {1..5}; do
    FILE="${APP_PKG_NAME}_${OS}_${ARCH}${EXT}"
    checksumURL="https://github.com/${ORG}/${REPO}/releases/latest/download/checksums.txt"
    resp=$(curl -Ls "${checksumURL}" -o /tmp/pct_checksums.txt --write-out "%{http_code}")
    respCode=$(echo ${resp} | tail -n 1)
    logDebug "GET ${checksumURL} | Resp: ${resp}"
    if [ ${respCode} -ne 200 ]; then
      echo "Fetching checksums.txt failed on attempt ${i}, retrying..."
      sleep 5
    else
      CHECKSUM=$(grep " ${FILE}" /tmp/pct_checksums.txt | cut -d ' ' -f 1)
      return 0
    fi
  done
  echo "Fetching checksums.txt failed after max retry attempts"
  exit 1
}

downloadLatestRelease() {
  destination="${HOME}/.puppetlabs/pct"

  [ -d ${destination} ] || mkdir -p ${destination} ]

  if [ "${noTel}" = "--no-telemetry" ]; then
      echo "Downloading and extracting ${APP_PKG_NAME} (TELEMETRY DISABLED VERSION) to ${destination}"
  else
      echo "Downloading and extracting ${APP_PKG_NAME} to ${destination}"
  fi

  downloadURL="https://github.com/${ORG}/${REPO}/releases/latest/download/${FILE}"

  for i in {1..5}; do
    resp=$(curl -Ls ${downloadURL} -o /tmp/${FILE} --write-out "%{http_code}")
    respCode=$(echo ${resp} | tail -n 1)
    logDebug "GET ${downloadURL} | Resp: ${resp}"
    if [ ${respCode} -ne 200 ]; then
      echo "Fetching PCT package failed on attempt ${i}, retrying..."
      sleep 5
    else
      downloadChecksumRaw=$(shasum -a 256 /tmp/${FILE} || sha256sum /tmp/${FILE})
      downloadChecksum=$(echo ${downloadChecksumRaw} | cut -d ' ' -f 1)
      logDebug "Checksum calc for ${FILE}:"
      logDebug " - Expect checksum: ${CHECKSUM}"
      logDebug " - Actual checksum: ${downloadChecksum}"
      if [ ${downloadChecksum} = ${CHECKSUM} ]; then
        logDebug "Extracting /tmp/${FILE} to ${destination}"
        tar -zxf "/tmp/${FILE}" -C ${destination}
        tarStatus=$(echo $?)
        logDebug "Removing /tmp/${FILE}"
        rm "/tmp/${FILE}"
        if [ ${tarStatus} -eq 0 ]; then
          echo "Remember to add the pct app to your path:"
          echo 'export PATH=$PATH:'${destination}
          exit 0
        else
          echo "Untar unsuccessful (status code: $?)"
          exit 1
        fi
      else
        echo "Checksum verification failed for ${FILE}"
        exit 1
      fi
      return 0
    fi
  done
}

getChecksums
downloadLatestRelease
