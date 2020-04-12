REPO="pixelandtonic/nitro"
VERSION="latest"
GITHUB_API_ENDPOINT="api.github.com"

if [ -z "$NITRO_TOKEN" ]; then
    echo "Missing \$NITRO_TOKEN"
    exit 1
fi

alias errcho='>&2 echo'

function gh_curl() {
  curl -sL -H "Authorization: token $NITRO_TOKEN" -H "Accept: application/vnd.github.v3.raw" $@
}
os=$(uname | tr '[:upper:]' '[:lower:]')
arch=$(uname -m)

FILE="nitro_$version"_"$os"_"$arch".tar.gz
if [ "$VERSION" = "latest" ]; then
  # Github should return the latest release first.
  PARSER=".[0].assets | map(select(.name == \"$FILE\"))[0].id"
else
  PARSER=". | map(select(.tag_name == \"$VERSION\"))[0].assets | map(select(.name == \"$FILE\"))[0].id"
fi

ASSET_ID=`gh_curl https://$GITHUB_API_ENDPOINT/repos/$REPO/releases | jq "$PARSER"`
if [ "$ASSET_ID" = "null" ]; then
  errcho "ERROR: version not found $VERSION"
  exit 1
fi

curl -sL --header "Authorization: token $NITRO_TOKEN" --header 'Accept: application/octet-stream' https://$NITRO_TOKEN:@$GITHUB_API_ENDPOINT/repos/$REPO/releases/assets/$ASSET_ID > $FILE