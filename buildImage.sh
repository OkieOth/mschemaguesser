#!/bin/bash

scriptPos=${0%/*}


imageBase=ghcr.io/okieoth/mschemaguesser
imageTag=`cat $scriptPos/cmd/schemaguesser/cmd/version.go | grep "const Version =" | sed -e 's-.*= "--' -e 's-".*--'`

imageName="$imageBase:$imageTag"

echo "I am going to create: $imageName"

pushd "$scriptPos" > /dev/null
        if docker build -f Dockerfile.release -t $imageName .
    then
        docker tag $imageName $imageBase
        echo -en "\033[1;34m  image created: $imageName, $imageBase \033[0m\n"
    else
        echo -en "\033[1;31m  error while create image: $imageName \033[0m\n"
    fi
popd > /dev/null