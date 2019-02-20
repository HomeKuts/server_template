#!/bin/bash

# init
USER="name"
WORK_DIR="/home/name/lalala"
SERVICE="cmd"




# start
cd $WORK_DIR
CMD="$WORK_DIR/$SERVICE"
PARAM="-config ./$SERVICE.ini"


# check if proccess already running
S=`pgrep -f $CMD`

if [[ -z $S ]]
then
    # sudo -u $USER $CMD $PARAM >> $WORK_DIR/service.log 2>&1 &
    $CMD $PARAM >> $WORK_DIR/service.log 2>&1 &
else
    echo "Process already running"
fi
