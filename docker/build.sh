#!/bin/sh

rootPath=$GOPATH
projectName=magicBlog
projectPath=$rootPath/src/github.com/muidea/$projectName
binPath=$rootPath/bin/$projectName
imageID=""
imageNamespace=muidea.ai/develop
imageVersion=latest
imageName=$imageNamespace/$(echo $projectName | tr '[A-Z]' '[a-z]')

function cleanUp()
{
    echo "cleanUp..."
    if [ -f log.txt ]; then
        rm -f log.txt
    fi

    if [ -f $projectName ]; then
        rm -f $projectName
    fi

    if [ -f $binPath ]; then
        rm -f $binPath
    fi

    if [ -f static ]; then
        rm -f static
    fi
}

function buildBin()
{
    echo "buildBin..."
    GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -tags netgo  -o $binPath github.com/muidea/magicBlog/cmd/magicBlog
    if [ $? -ne 0 ]; then
        echo "buildBin failed."
        exit 1
    else
        echo "buildBin success."
    fi
}

function prepareFile()
{
    echo "prepareFile..."
    if [ ! -f $binPath ]; then
        buildBin
        if [ $? -ne 0 ]; then
            exit 1
        fi
    fi

    cp $binPath ./
    if [ $? -ne 0 ]; then
        echo "prepareFile failed."
        exit 1
    else
        echo "prepareFile success."
    fi

    prePath=$(pwd)
    cd $projectPath
    tar -caf $prePath/static static
    if [ $? -ne 0 ]; then
        echo "prepare file failed, tar static files failed exception."
        exit 1
    else
        echo "tar static files filed success."
    fi
    cd $prePath
}

function checkImage()
{
    echo "checkImage..."
    docker images | grep $1 | grep $2 > log.txt
    if [ $? -eq 0 ]; then
        imageID=$(tail -1 log.txt|awk '{print $3}')
    fi
}

function buildImage()
{
    echo "buildImage..."
    docker build . > log.txt
    if [ $? -eq 0 ]; then
        echo "buildImage success."
    else
        echo "buildImage failed."
        exit 1
    fi

    imageID=$(tail -1 log.txt|awk '{print $3}')
}


function tagImage()
{
    echo "tagImage image..."
    docker tag $1 $2
    if [ $? -eq 0 ]; then
        echo "tagImage success."
    else
        echo "tagImage failed."
        exit 1
    fi
}

function rmiImage()
{
    echo "rmiImage..."
    docker rmi $1:$2
    if [ $? -eq 0 ]; then
        echo "rmiImage success."
    else
        echo "rmiImage failed."
        exit 1
    fi
}

function all()
{
    echo "build magicBlog docker image"

    curPath=$(pwd)

    cd $projectPath/docker

    cleanUp

    prepareFile

    checkImage $imageName $imageVersion
    if [ $imageID ]; then
        rmiImage $imageName $imageVersion
    fi

    buildImage

    tagImage $imageID $imageName:$imageVersion

    cleanUp

    cd $curPath
}

function build()
{
    checkImage $imageName $imageVersion
    if [ $imageID ]; then
        rmiImage $imageName $imageVersion
    fi

    buildImage

    tagImage $imageID $imageName:$imageVersion    
}

action='all'
if [ $1 ]; then
    action=$1
fi

if [ $action == 'prepare' ]; then
    prepareFile
elif [ $action == 'clean' ]; then
    cleanUp
elif [ $action == 'build' ]; then
    build
else
    all
fi
