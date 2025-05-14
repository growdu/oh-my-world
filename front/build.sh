#!/bin/bash
yarn run build-only
rm -rf nginx/dist_bak
if [ -e "nginx/dist" ]
then
	cp nginx/dist nginx/dist_bak
fi
cp dist nginx/ -r
