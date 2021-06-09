#!/usr/bin/env bash

# Run database migrations.
/app/urlshortener migrate
if [[ $? -ne 0 ]] ; then
    exit 1
fi

# Run urlshortener server.
/app/urlshortener server
