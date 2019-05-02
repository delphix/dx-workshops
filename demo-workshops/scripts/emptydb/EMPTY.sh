#!/bin/sh
#
# Copyright (c) 2019 by Delphix. All rights reserved.
#
source ~/.bash_profile
DBPASSWORD=${1}
echo "Oracle base is ${ORACLE_BASE}"
echo "Oracle home is ${ORACLE_HOME}"
OLD_UMASK=`umask`
umask 0027
mkdir -p ${ORACLE_BASE}/admin/EMPTY/adump
mkdir -p ${ORACLE_BASE}/admin/EMPTY/dpdump
mkdir -p ${ORACLE_BASE}/admin/EMPTY/pfile
mkdir -p ${ORACLE_BASE}/cfgtoollogs/dbca/EMPTY
mkdir -p ${ORACLE_BASE}/oradata/EMPTY
mkdir -p ${ORACLE_HOME}/dbs
umask ${OLD_UMASK}
ORACLE_SID=EMPTY; export ORACLE_SID
PATH=$ORACLE_HOME/bin:$PATH; export PATH
echo "EMPTY:${ORACLE_HOME}:Y" >> /etc/oratab
${ORACLE_HOME}/bin/sqlplus /nolog @/home/oracle/emptydb/EMPTY.sql <<EOM
${DBPASSWORD}
${DBPASSWORD}
${DBPASSWORD}
EOM
/home/oracle/emptydb/createDelphixDBUser.sh <<EOM
EMPTY
delphixdb
${DBPASSWORD}
y
EOM
/${ORACLE_HOME}/bin/sqlplus / as sysdba @/home/oracle/emptydb/target_grants.sql

/home/oracle/emptydb/ready.sh