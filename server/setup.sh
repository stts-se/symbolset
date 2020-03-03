#!/bin/bash 

set -e

CMD=`basename $0`
SCRIPTDIR=`dirname $0`

if [ $# -ne 2 ]; then
    echo "USAGE: bash $CMD <LEXDATA-DIR> <DEST-DIR>" >&2
    exit 1
fi

LEXDATA=$1
DESTDIR=$2

mkdir -p $DESTDIR

echo -n "Copying demo files ..." >&2
cp $SCRIPTDIR/demo_files/* $DESTDIR
echo " done" >&2

echo -n "Copying symbol sets ..." >&2
cp $LEXDATA/*/*/*.sym $DESTDIR
echo " done" >&2

echo -n "Copying converters ..." >&2
cp $LEXDATA/converters/*.cnv $DESTDIR
echo " done" >&2

echo -n "Copying mappers ..." >&2
cat $LEXDATA/mappers.txt >> $DESTDIR/mappers.txt
echo " done" >&2
