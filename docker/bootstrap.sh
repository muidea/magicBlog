#!/bin/bash

EXTRA_ARGS=$EXTRA_ARGS
SETUP_ARGS=$SETUP_ARGS
if [ $LISTENPORT ]; then
    EXTRA_ARGS=$EXTRA_ARGS' -ListenPort='$LISTENPORT
fi

if [ $CENTERSERVER ]; then
    EXTRA_ARGS=$EXTRA_ARGS' -CenterSvr='$CENTERSERVER
    SETUP_ARGS=$SETUP_ARGS' -CenterSvr='$CENTERSERVER
fi

if [ $ENDPOINTNAME ]; then
    EXTRA_ARGS=$EXTRA_ARGS' -EndpointName='$ENDPOINTNAME
    SETUP_ARGS=$SETUP_ARGS' -EndpointName='$ENDPOINTNAME
fi

if [ $ENDPOINTID ]; then
    EXTRA_ARGS=$EXTRA_ARGS' -EndpointID='$ENDPOINTID
    SETUP_ARGS=$SETUP_ARGS' -EndpointID='$ENDPOINTID
fi

if [ $AUTHTOKEN ]; then
    EXTRA_ARGS=$EXTRA_ARGS' -AuthToken='$AUTHTOKEN
    SETUP_ARGS=$SETUP_ARGS' -AuthToken='$AUTHTOKEN
fi

echo $EXTRA_ARGS

/var/app/wait-for-it.sh $CENTERSERVER -- echo "centerServer is ready."

/var/app/setupTool $SETUP_ARGS "$@"

/var/app/magicBlog $EXTRA_ARGS "$@"
