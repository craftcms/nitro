#!/bin/bash

export GH_ORG=craftcms
export SUCCESS_CMD="nitro"
export BINLOCATION="/usr/local/bin"
export TEMP_FOLDER="temp_nitro_extract"

version=$(curl -s https://api.github.com/repos/craftcms/nitro/releases/latest | grep -i tag_name | sed 's/\(\"tag_name\": \"\(.*\)\",\)/\2/' | tr -d '[:space:]')

if [ ! "$version" ]; then
  echo "There was a problem trying to automatically install Nitro. You can try to install manually:"
  echo
  echo "1. Open your web browser and go to https://github.com/craftcms/nitro/releases"
  echo "2. Download the latest release for your platform and unzip it."
  echo "3. Run 'chmod +x ./nitro' on the unzipped 'nitro' executable."
  echo "4. Run 'mv ./nitro $BINLOCATION'"
  exit 1
fi

hasCurl() {
  result=$(command -v curl)
  if [ "$?" = "1" ]; then
    echo "You need curl to use this script."
    exit 1
  fi
}

hasMultipass() {
  result=$(command -v multipass)
  if [ "$?" = "1" ]; then
    echo "You need Multipass to use Nitro. Please install it for your platform https://multipass.run/"
    exit 1
  fi
}

checkHash() {
  sha_cmd="sha256sum"
  fileName=nitro_$2_checksums.txt
  filePath=$(pwd)/$TEMP_FOLDER/$fileName
  checksumUrl=https://github.com/craftcms/nitro/releases/download/$version/$fileName
  targetFile=$3/$fileName

  if [ ! -x "$(command -v $sha_cmd)" ]; then
    shaCmd="shasum -a 256"
  fi

  if [ -x "$(command -v $shaCmd)" ]; then

    # download the checksum file.
    (curl -sLS "$checksumUrl" --output "$targetFile")

    # Run the sha command against the zip and grab the hash from the first segment.
    zipHash="$($shaCmd $1 | cut -d' ' -f1 | tr -d '[:space:]')"

    # See if the has we calculated matches a result in the checksum file.
    checkResultFileName=$(sed -n "s/^$zipHash  //p" "$filePath")

    # don't need this anymore
    rm "$filePath"

    # Make sure the file names match up.
    if [ "$4" != "$checkResultFileName" ]; then
      rm "$1"
      echo "Checksums do not match. Exiting."
      exit 1
    fi
  fi
}

getNitro() {
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
    suffix="_linux"

    case $arch in
    "aarch64")
      suffix="_linux"
      ;;
    esac
    ;;
  esac

  if [ ! -d $(pwd)/$TEMP_FOLDER ]; then
    mkdir $(pwd)/$TEMP_FOLDER
  fi

  targetTempFolder="$(pwd)/$TEMP_FOLDER"

  fileName=nitro"$suffix"_x86_64.tar.gz
  packageUrl=https://github.com/craftcms/nitro/releases/download/$version/"$fileName"
  targetZipFile="$targetTempFolder"/$fileName

  echo "Downloading package $packageUrl to $targetZipFile"
  echo
  curl -sSL "$packageUrl" --output "$targetZipFile"

  if [ "$?" = "0" ]; then

    # unzip
    tar xvzf "$targetZipFile" -C "$targetTempFolder"

    # verify
    checkHash "$targetZipFile" "$version" "$targetTempFolder" "$fileName"

    mv "$targetTempFolder"/nitro ./nitro
    chmod +x ./nitro

    echo
    echo "Download complete."

    if [ ! -w "$BINLOCATION" ]; then
      echo
      echo "============================================================"
      echo "  The script was run as a user who is unable to write"
      echo "  to $BINLOCATION. To complete the installation the"
      echo "  following commands may need to be run manually."
      echo "============================================================"
      echo
      echo "  sudo cp ./nitro $BINLOCATION/nitro"
      echo "  nitro"
      echo
    else
      echo
      echo "Running with sufficient permissions to attempt to move the nitro executable to $BINLOCATION"
      echo

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
        echo "Nitro $version has been installed to $BINLOCATION"
      fi

      if [ -e "$targetTempFolder" ]; then
        rm -rf "$targetTempFolder"
        echo
      fi

      ${SUCCESS_CMD}
      init
    fi
  fi
}

init() {
  echo
  read -p "Initialize the primary machine now? " -n 1 -r
  echo
  if [[ $REPLY =~ ^[Yy]$ ]]; then
    nitro init
  fi
}

hasCurl
hasMultipass
getNitro
