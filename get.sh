#!/bin/bash

export SYMLINK_NAME="nitro"
export GH_ORG=pixelandtonic
export REPO=nitro
export SUCCESS_CMD="$REPO version"
export BINLOCATION="/usr/local/bin"

version=$(curl -s https://$NITRO_TOKEN@api.github.com/repos/$GH_ORG/$REPO/releases/latest | grep -i tag_name | sed 's/\(\"tag_name\": \"\(.*\)\",\)/\2/' | tr -d '[:space:]')

if [ ! "$version" ]; then
  echo "There was a problem trying to automatically install $REPO. You can try to install manually:"
  echo ""
  echo "1. Open your web browser and go to https://github.com/$GH_ORG/$REPO/releases"
  echo "2. Download the latest release for your platform. Call it '$REPO'."
  echo "3. chmod +x ./$REPO"
  echo "4. mv ./$REPO $BINLOCATION"
  if [ -n "$SYMLINK_NAME" ]; then
    echo "5. ln -sf $BINLOCATION/$REPO /usr/local/bin/$SYMLINK_NAME"
  fi
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
  checksumUrl=https://github.com/$GH_ORG/$REPO/releases/download/$version/checksums.txt

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
      suffix="-darwin"
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
          suffix="-arm64"
          ;;
      esac

      case $arch in
        "armv6l" | "armv7l")
          suffix="-armhf"
          ;;
      esac
    ;;
  esac

  targetFile="/tmp/$REPO$suffix"

  if [ "$userid" != "0" ]; then
    targetFile="$(pwd)/$REPO$suffix"
  fi

  if [ -e $targetFile ]; then
    rm "$targetFile"
  fi

  packageUrl=https://github.com/$GH_ORG/$REPO/releases/download/$version/$REPO$suffix
  echo "Downloading package $packageUrl as $targetFile"

  curl -sSL "$packageUrl" --output "$targetFile"

  if [ "$?" = "0" ]; then
    checkHash
    chmod +x "$targetFile"
    echo "Download complete."

    if [ ! -w "$BINLOCATION" ]; then
      echo
      echo "============================================================"
      echo "  The script was run as a user who is unable to write"
      echo "  to $BINLOCATION. To complete the installation the"
      echo "  following commands may need to be run manually."
      echo "============================================================"
      echo
      echo "  sudo cp $REPO$suffix $BINLOCATION/$REPO"

      if [ -n "$SYMLINK_NAME" ]; then
          echo "  sudo ln -sf $BINLOCATION/$REPO $BINLOCATION/$SYMLINK_NAME"
      fi

      echo
    else
      echo
      echo "Running with sufficient permissions to attempt to move $REPO to $BINLOCATION"

      if [ ! -w "$BINLOCATION/$REPO" ] && [ -f "$BINLOCATION/$REPO" ]; then
        echo
        echo "================================================================"
        echo "  $BINLOCATION/$REPO already exists and is not writeable"
        echo "  by the current user.  Please adjust the binary ownership"
        echo "  or run sh/bash with sudo."
        echo "================================================================"
        echo
        exit 1
      fi

      mv "$targetFile" "$BINLOCATION"/$REPO

      if [ "$?" = "0" ]; then
        echo "New version of $REPO installed to $BINLOCATION"
      fi

      if [ -e "$targetFile" ]; then
        rm "$targetFile"
      fi

      if [ -n "$SYMLINK_NAME" ]; then
        if [ ! -L "$BINLOCATION"/$SYMLINK_NAME ]; then
          ln -s "$BINLOCATION"/$REPO "$BINLOCATION"/$SYMLINK_NAME
          echo "Creating symlink '$SYMLINK_NAME' for '$REPO'."
        fi
      fi

      ${SUCCESS_CMD}
    fi
  fi
}

hasCurl
getNitro