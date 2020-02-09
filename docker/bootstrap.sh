#!/bin/bash

EXTRA_ARGS=$EXTRA_ARGS
if [ $LISTENPORT ]; then
    EXTRA_ARGS=$EXTRA_ARGS' -ListenPort='$LISTENPORT
fi

if [ $ENDPOINTNAME ]; then
    EXTRA_ARGS=$EXTRA_ARGS' -EndpointName='$ENDPOINTNAME
fi

if [ $CASSERVER ]; then
    EXTRA_ARGS=$EXTRA_ARGS' -CasService='$CASSERVER
fi

if [ $USERSERVER ]; then
    EXTRA_ARGS=$EXTRA_ARGS' -UserService='$USERSERVER
fi

if [ $FILESERVER ]; then
    EXTRA_ARGS=$EXTRA_ARGS' -FileService='$FILESERVER
fi

if [ $BATISSERVER ]; then
    EXTRA_ARGS=$EXTRA_ARGS' -BatisService='$BATISSERVER
fi

echo $EXTRA_ARGS

/var/app/magicBlog $EXTRA_ARGS "$@"
