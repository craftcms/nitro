#!/bin/bash

export TEMP_FOLDER="temp_nitro_extract"
export FINAL_DIR_LOCATION="/usr/local/bin"
export DOWNLOAD_SUFFIX=""
export DOWNLOAD_ZIP_EXTENSION=".tar.gz"
export IS_WINDOWS=false
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
    IS_WINDOWS=true
    DOWNLOAD_SUFFIX="_windows"
    DOWNLOAD_ZIP_EXTENSION=".zip"
    FINAL_DIR_LOCATION="$HOME/Nitro"
    EXECUTABLE_NAME="nitro.exe"
    ;;

  esac
}

checkPlatform

version=$(curl -s https://api.github.com/repos/craftcms/nitro/releases/latest | grep -i tag_name | sed 's/\(\"tag_name\": \"\(.*\)\",\)/\2/' | tr -d '[:space:]')

if [ ! "$version" ]; then
  echo "There was a problem trying to automatically install Nitro. You can try to install manually:"
  echo

  if [ "$IS_WINDOWS" = true ]; then
    echo "1. Open your web browser and go to https://github.com/craftcms/nitro/releases"
    echo "2. Download \"nitro_windows_x86_64.zip\" for latest release of Nitro and unzip it."
    echo "3. Make sure $FINAL_DIR_LOCATION exists, then copy $EXECUTABLE_NAME from the unzipped folder into it."
    echo "4. Open Git Bash and run the following commands:"
    echo
    echo "       export PATH=$FINAL_DIR_LOCATION:\$PATH"
    echo "       $EXECUTABLE_NAME"
  else
    echo "1. Open your web browser and go to https://github.com/craftcms/nitro/releases"
    echo "2. Download the latest release for your platform and unzip it."
    echo "3. Run 'chmod +x ./$EXECUTABLE_NAME' on the unzipped \"$EXECUTABLE_NAME\" executable."
    echo "4. Run 'mv ./$EXECUTABLE_NAME $FINAL_DIR_LOCATION'"
    echo "5. Run 'nitro init' to create your first machine."
  fi

  exit 1
fi

function hasCurl {
  if [ "$IS_WINDOWS" = true ]; then
    result=$(where curl)
    if [[ "$?" == *"Could not find files"* ]]; then
      echo "You need curl to use Nitro."
      exit 1
    fi
  else
    result=$(command -v curl)
    if [ "$?" = "1" ]; then
      echo "You need curl to use Nitro."
      exit 1
    fi
  fi
}

function hasMultipass {
  if [ "$IS_WINDOWS" = true ]; then
    result=$(where multipass)
    if [[ "$?" == *"Could not find files"* ]]; then
      echo "You need Multipass to use Nitro. Please install it for your platform https://multipass.run/"
      exit 1
    fi
  else
    result=$(command -v multipass)
    if [ "$?" = "1" ]; then
      echo "You need Multipass to use Nitro. Please install it for your platform https://multipass.run/"
      exit 1
    fi
  fi
}

function checkHash {
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
      echo "It looks like there was an incomplete download. Please try again."
      exit 1
    fi
  fi
}

function getNitro {
  # if it's Windows, make sure the final destination exists
  if [ "$IS_WINDOWS" = true ] && [ ! -d "$FINAL_DIR_LOCATION" ]; then
    mkdir -p "$FINAL_DIR_LOCATION"
  fi

  # create our temp folder
  if [ ! -d $(pwd)/$TEMP_FOLDER ]; then
    mkdir $(pwd)/$TEMP_FOLDER
  fi

  targetTempFolder="$(pwd)/$TEMP_FOLDER"

  fileName=nitro"$DOWNLOAD_SUFFIX"_x86_64"$DOWNLOAD_ZIP_EXTENSION"
  packageUrl=https://github.com/craftcms/nitro/releases/download/$version/"$fileName"
  targetZipFile="$targetTempFolder"/$fileName

  echo "Downloading package $packageUrl to $targetZipFile"
  echo
  curl -sSL "$packageUrl" --output "$targetZipFile"

  if [ "$?" = "0" ]; then

    # unzip
    if [ "$IS_WINDOWS" = true ]; then
      unzip "$targetZipFile" -d "$targetTempFolder"
    else
      tar xvzf "$targetZipFile" -C "$targetTempFolder"
    fi

    # verify
    checkHash "$targetZipFile" "$version" "$targetTempFolder" "$fileName"

    # move executable up a level
    mv "$targetTempFolder"/"$EXECUTABLE_NAME" ./"$EXECUTABLE_NAME"

    # make it executable
    if [ "$IS_WINDOWS" = false ]; then
      chmod +x ./$EXECUTABLE_NAME
    fi

    echo
    echo "Download complete."

    if [ ! -w "$FINAL_DIR_LOCATION" ]; then
      if [ "$IS_WINDOWS" = true ]; then
        echo
        echo "============================================================"
        echo "  The script was run as a user who is unable to write"
        echo "  to $FINAL_DIR_LOCATION. To complete the installation make"
        echo "  sure $FINAL_DIR_LOCATION exists, then copy $EXECUTABLE_NAME"
        echo "  from the current folder into it, then run the following:"
        echo "============================================================"
        echo
        echo "  export PATH=$FINAL_DIR_LOCATION:\$PATH"
        echo "  $EXECUTABLE_NAME init"
        echo
      else
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
      fi
    else
      echo
      echo "Running with sufficient permissions to attempt to move $EXECUTABLE_NAME to $FINAL_DIR_LOCATION"
      echo

      if [ ! -w "$FINAL_DIR_LOCATION/$EXECUTABLE_NAME" ] && [ -f "$FINAL_DIR_LOCATION/$EXECUTABLE_NAME" ]; then
        if [ "$IS_WINDOWS" = true ]; then
          echo
          echo "============================================================"
          echo "  $FINAL_DIR_LOCATION/$EXECUTABLE_NAME alraedy exists and is"
          echo "  not writable by this installer. To complete the installation make"
          echo "  sure $FINAL_DIR_LOCATION exists, then copy $EXECUTABLE_NAME"
          echo "  from the current folder into it, then run the following:"
          echo "============================================================"
          echo
          echo "  export PATH=$FINAL_DIR_LOCATION:\$PATH"
          echo "  $EXECUTABLE_NAME init"
          echo
          exit 1
        else
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
promptInstall
