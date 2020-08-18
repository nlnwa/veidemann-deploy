#!/usr/bin/env bash

# Pipe kubernetes manifests into this script and get a JSON object
# with image name and version as output
#
# Example usage:
# kustomize build | ./versions.sh | jq

# overall version
tag=$(git describe --tag --always --dirty)

# read standard input into variable
input=$(</dev/stdin)

# grep for images in kubernetes yaml syntax
images=$(echo "$input" | grep -e '^[[:blank:]]*image: [^ ]*' | awk '{print $2}')

# count number of images to be able to not append last comma in json output
count=$(echo "$images" | wc -l)

# start json object
VERSIONS_JSON+="{"

# append overall version
VERSIONS_JSON+="\"veidemann\": \"$tag\","

# append image versions
for version in $images; do
  ((count--))
  # split single version into parts delimited by colon
  IFS=':' read -a parts <<<${version}
  # append image version
  VERSIONS_JSON+="\"${parts[0]}\": \"${parts[1]}\""
  if ((count > 0)); then VERSIONS_JSON+=","; fi
done

# end json object
VERSIONS_JSON+="}"

# output images as json object
echo "$VERSIONS_JSON"
