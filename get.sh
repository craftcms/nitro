#!/bin/bash

export GH_ORG=pixelandtonic
export SUCCESS_CMD="nitro"
export BINLOCATION="/usr/local/bin"

version=$(curl -s https://api.github.com/repos/pixelandtonic/nitro/releases/latest | grep -i tag_name | sed 's/\(\"tag_name\": \"\(.*\)\",\)/\2/' | tr -d '[:space:]')

if [ ! "$version" ]; then
  echo "There was a problem trying to automatically install nitro. You can try to install manually:"
  echo ""
  echo "1. Open your web browser and go to https://github.com/pixelandtonic/nitro/releases"
  echo "2. Download the latest release for your platform. Call it 'nitro'."
  echo "3. chmod +x ./nitro"
  echo "4. mv ./nitro $BINLOCATION"
  exit 1
fi

hasCurl() {
  $(which curl)
  if [ "$?" = "1" ]; then
    echo "You need curl to use this script."
    exit 1
  fi
}

checkHash () {
  sha_cmd="sha256sum"
  checksumUrl=https://github.com/pixelandtonic/nitro/releases/download/$version/checksums.txt

  if [ ! -x "$(command -v $sha_cmd)" ]; then
    shaCmd="shasum -a 256"
  fi

  if [ -x "$(command -v $shaCmd)" ]; then

    targetFileDir=${targetFile%/*}

    (cd "$targetFileDir" && curl -sSL "$packageUrl")

    echo "cd $targetFileDir && curl -sSL $packageUrl.sha256|$shaCmd -c"

    (cd "$targetFileDir" && curl -sSL "$packageUrl"|$shaCmd -c >/dev/null)

    if [ "$?" != "0" ]; then
      rm "$targetFile"
      echo "Binary checksum didn't match. Exiting"
      exit 1
    fi
  fi
}

getNitro () {
  uname=$(uname)
  userid=$(id -u)

  suffix=""
  case $uname in

    "Darwin")
      suffix="_darwin"
      ;;

    "MINGW"*)
      suffix=".exe"
      BINLOCATION="$HOME/bin"
      mkdir -p "$BINLOCATION"
      ;;

    "Linux")
      arch=$(uname -m)
      echo "$arch"

      case $arch in
        "aarch64")
          suffix="_arm64"
          ;;
      esac

      case $arch in
        "armv6l" | "armv7l")
          suffix="_armhf"
          ;;
      esac
    ;;
  esac

  targetTempFolder="/tmp"

  if [ "$userid" != "0" ]; then
    targetTempFolder="$(pwd)"
  fi

  fileName=nitro_"$version""$suffix"_x86_64.tar.gz
  packageUrl=https://github.com/pixelandtonic/nitro/releases/download/$version/"$fileName"
  targetZipFile="$targetTempFolder"/$fileName

  echo "Downloading package $packageUrl to $targetZipFile"

  curl -sSL "$packageUrl" --output "$targetZipFile"

  if [ "$?" = "0" ]; then

    #unzip
    tar xvzf "$targetZipFile"

    # TODO add checkHash
    # checkHash
    chmod +x ./nitro
    echo "Download complete."

    if [ ! -w "$BINLOCATION" ]; then
      echo
      echo "============================================================"
      echo "  The script was run as a user who is unable to write"
      echo "  to $BINLOCATION. To complete the installation the"
      echo "  following commands may need to be run manually."
      echo "============================================================"
      echo
      echo "  sudo cp nitro$suffix $BINLOCATION/nitro"
      echo
    else
      echo
      echo "Running with sufficient permissions to attempt to move nitro to $BINLOCATION"

      if [ ! -w "$BINLOCATION/nitro" ] && [ -f "$BINLOCATION/nitro" ]; then
        echo
        echo "================================================================"
        echo "  $BINLOCATION/nitro already exists and is not writeable"
        echo "  by the current user.  Please adjust the binary ownership"
        echo "  or run sh/bash with sudo."
        echo "================================================================"
        echo
        exit 1
      fi

      mv ./nitro "$BINLOCATION"/nitro

      if [ "$?" = "0" ]; then
        echo "New version of nitro installed to $BINLOCATION"
        echo
      fi

      if [ -e "$targetNitroFile" ]; then
        rm "$targetNitroFile"
      fi

      ${SUCCESS_CMD}
    fi
  fi
}

hasCurl
getNitro