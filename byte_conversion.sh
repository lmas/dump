#!/bin/sh

# http://www.linuxquestions.org/questions/linux-general-1/awk-to-convert-bytes-to-human-number-909214/#post4504048

echo "$1" | awk '{ split( "B KB MB GB TB" , v ); s=1; while( $1>1024 ){ $1/=1024; s++ } print int($1) v[s] }'
