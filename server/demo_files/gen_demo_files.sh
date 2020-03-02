#!/bin/bash 

CMD=`basename $0`
SCRIPTDIR=`dirname $0`

if [ $# -ne 2 ]; then
    echo "USAGE: bash $CMD <LEXDATA-DIR> <DEST-DIR>" >&2
    exit 1
fi

LEXDATA=$1
DESTDIR=$2

mkdir -p $DESTDIR

echo "Copying symbol sets" >&2
cp $LEXDATA/*/*/*.sym $DESTDIR

echo "Copying converters" >&2
cp $LEXDATA/converters/*.cnv $DESTDIR

echo "Copying mappers" >&2
cp $LEXDATA/mappers.txt $DESTDIR
