#!/bin/bash

EXTRA_ARGS=$EXTRA_ARGS
if [ $LISTENPORT ]; then
    EXTRA_ARGS=$EXTRA_ARGS' -ListenPort='$LISTENPORT
fi

if [ $CENTER_SERVER ]; then
    EXTRA_ARGS=$EXTRA_ARGS' -CenterSvr='$CENTER_SERVER
fi

if [ $NAME ]; then
    EXTRA_ARGS=$EXTRA_ARGS' -BlogName='$NAME
fi

if [ $ENDPOINTID ]; then
    EXTRA_ARGS=$EXTRA_ARGS' -EndpointID='$ENDPOINTID
fi

if [ $AUTHTOKEN ]; then
    EXTRA_ARGS=$EXTRA_ARGS' -AuthToken='$AUTHTOKEN
fi

echo $EXTRA_ARGS

/var/app/magicBlog $EXTRA_ARGS "$@"
