#!/bin/sh
source ~/.bash_profile
echo "Oracle home is $ORACLE_HOME"
ORACLE_SID=EMPTY; export ORACLE_SID
PATH=$ORACLE_HOME/bin:$PATH; export PATH

CONNECT="/ as sysdba"
until [[ ${READY} = "OPEN" ]]; do
    echo "Testing ${ORACLE_HOME}"
    READY=$(sqlplus "${CONNECT}" @/home/oracle/emptydb/ready.sql |grep OPEN | awk '{print $1}')
    sleep 2
done
