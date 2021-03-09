#!/bin/bash

export TEMP_FOLDER="temp_nitro_extract"
export FINAL_DIR_LOCATION="/usr/local/bin"
export DOWNLOAD_SUFFIX=""
export DOWNLOAD_ARCH=""
export DOWNLOAD_ZIP_EXTENSION=".tar.gz"
export EXECUTABLE_NAME="nitro"

function checkPlatform {
  uname=$(uname)

  case $uname in

  "Darwin")
    DOWNLOAD_SUFFIX="_darwin"
    ;;

  "Linux")
    DOWNLOAD_SUFFIX="_linux"
    ;;

  *)
    ;;

  esac

  arch=$(uname -m)

  case $arch in

  "x86_64")
    DOWNLOAD_ARCH="_x86_64"
    ;;

  "aarch64")
    DOWNLOAD_ARCH="_aarch64"
    ;;

  "arm64")
    DOWNLOAD_ARCH="_arm64"
    ;;

  *)
    ;;

  esac
}

checkPlatform

version=$(curl -s https://api.github.com/repos/craftcms/nitro/releases | grep -i -m 1 tag_name | head -1 | sed 's/\("tag_name": "\(.*\)",\)/\2/' | tr -d '[:space:]')

if [ ! "$version" ]; then
  echo "There was a problem trying to automatically install Nitro. You can try to install manually:"
  echo

  echo "1. Open your web browser and go to https://github.com/craftcms/nitro/releases"
  echo "2. Download the latest release for your platform and unzip it."
  echo "3. Run 'chmod +x ./$EXECUTABLE_NAME' on the unzipped \"$EXECUTABLE_NAME\" executable."
  echo "4. Run 'mv ./$EXECUTABLE_NAME $FINAL_DIR_LOCATION'"
  echo "5. Run 'nitro init' to create your first machine."

  exit 1
fi

function hasCurl {
  result=$(command -v curl)
  if [ "$?" = "1" ]; then
    echo "You need curl to install Nitro."
    exit 1
  fi
}

function hasDocker {
  result=$(command -v docker)
  if [ "$?" = "1" ]; then
    echo "You need Docker Desktop to use Nitro. Please install it for your platform at https://www.docker.com/products/docker-desktop"
    exit 1
  fi
}

function checkHash {
  sha_cmd="sha256sum"
  fileName=nitro_$2_checksums.txt
  filePath="$(pwd)/$TEMP_FOLDER/$fileName"
  checksumUrl=https://github.com/craftcms/nitro/releases/download/$version/$fileName
  targetFile=$3/$fileName

  if [ ! -x "$(command -v $sha_cmd)" ]; then
    shaCmd="shasum -a 256"
  fi

  if [ -x "$(command -v $shaCmd)" ]; then

    # download the checksum file.
    (curl -sLS "$checksumUrl" --output "$targetFile")

    # Run the sha command against the zip and grab the hash from the first segment.
    zipHash="$($shaCmd "$1" | cut -d' ' -f1 | tr -d '[:space:]')"

    # See if the has we calculated matches a result in the checksum file.
    checkResultFileName=$(sed -n "s/^$zipHash  //p" "$filePath")

    # don't need this anymore
    rm "$filePath"

    # Make sure the file names match up.
    if [ "$4" != "$checkResultFileName" ]; then
      rm "$1"
      echo "It looks like there was an incomplete download. Please try again."
      exit 1
    fi
  fi
}

function getNitro {
  targetTempFolder="$(pwd)/$TEMP_FOLDER"

  # create our temp folder
  if [ ! -d "$targetTempFolder" ]; then
    mkdir "$targetTempFolder"
  fi

  fileName=nitro$DOWNLOAD_SUFFIX$DOWNLOAD_ARCH$DOWNLOAD_ZIP_EXTENSION
  packageUrl=https://github.com/craftcms/nitro/releases/download/$version/"$fileName"
  targetZipFile="$targetTempFolder/$fileName"

  echo "Downloading package $packageUrl to $targetZipFile"
  echo
  curl -sSL "$packageUrl" --output "$targetZipFile"

  if [ "$?" = "0" ]; then

    # unzip
    tar xvzf "$targetZipFile" -C "$targetTempFolder"

    # verify
    checkHash "$targetZipFile" "$version" "$targetTempFolder" "$fileName"

    # move executable up a level
    mv "$targetTempFolder"/"$EXECUTABLE_NAME" ./"$EXECUTABLE_NAME"

    # make it executable
    chmod +x ./$EXECUTABLE_NAME

    echo
    echo "Download complete."

    if [ ! -w "$FINAL_DIR_LOCATION" ]; then
      echo
      echo "============================================================"
      echo "  The script was run as a user who is unable to write"
      echo "  to $FINAL_DIR_LOCATION. To complete the installation the"
      echo "  following commands may need to be run manually:"
      echo "============================================================"
      echo
      echo "  sudo cp ./$EXECUTABLE_NAME $FINAL_DIR_LOCATION/$EXECUTABLE_NAME"
      echo "  $EXECUTABLE_NAME init"
      echo
      exit 1
    else
      echo
      echo "Running with sufficient permissions to attempt to move $EXECUTABLE_NAME to $FINAL_DIR_LOCATION"
      echo

      if [ ! -w "$FINAL_DIR_LOCATION/$EXECUTABLE_NAME" ] && [ -f "$FINAL_DIR_LOCATION/$EXECUTABLE_NAME" ]; then
        echo
        echo "============================================================"
        echo "  $FINAL_DIR_LOCATION/$EXECUTABLE_NAME already exists and is"
        echo "  not writeable by this installer. To complete the installation the"
        echo "  following commands may need to be run manually:"
        echo "============================================================"
        echo
        echo "  sudo cp ./$EXECUTABLE_NAME $FINAL_DIR_LOCATION/$EXECUTABLE_NAME"
        echo "  $EXECUTABLE_NAME init"
        echo
        exit 1
      fi

      # move to final location
      mv ./$EXECUTABLE_NAME "$FINAL_DIR_LOCATION"/$EXECUTABLE_NAME

      if [ "$?" = "0" ]; then
        echo "Nitro $version has been installed to $FINAL_DIR_LOCATION"
      fi

      if [ -e "$targetTempFolder" ] && [ -w "$targetTempFolder" ]; then
        rm -rf "$targetTempFolder"
        echo
      fi
    fi
  fi
}

function promptInstall {
    nitro init
}

hasCurl
hasDocker
getNitro

if [ "$1" != "update" ]; then
  promptInstall
fi
