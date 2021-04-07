#!/bin/sh

EXTRA_ARGS=$EXTRA_ARGS
if [ $LISTENPORT ]; then
    EXTRA_ARGS=$EXTRA_ARGS' -ListenPort='$LISTENPORT
fi

if [ $ENDPOINTNAME ]; then
    EXTRA_ARGS=$EXTRA_ARGS' -EndpointName='$ENDPOINTNAME
fi

if [ $BATISSERVICE ]; then
    EXTRA_ARGS=$EXTRA_ARGS' -BatisService='$BATISSERVICE
fi

if [ $CMSSERVICE ]; then
    EXTRA_ARGS=$EXTRA_ARGS' -CMSService='$CMSSERVICE
fi

if [ $CMSCATALOG ]; then
    EXTRA_ARGS=$EXTRA_ARGS' -CMSCatalog='$CMSCATALOG
fi

if [ $IDENTITYID ]; then
    EXTRA_ARGS=$EXTRA_ARGS' -IdentityID='$IDENTITYID
fi

if [ $AUTHTOKEN ]; then
    EXTRA_ARGS=$EXTRA_ARGS' -AuthToken='$AUTHTOKEN
fi


echo $EXTRA_ARGS

/var/app/magicBlog $EXTRA_ARGS "$@"
