#!/bin/sh
#
# Copyright (c) 2019 by Delphix. All rights reserved.
#
source ~/.bash_profile
ORACLE_SID=EMPTY; export ORACLE_SID
echo "Oracle base is ${ORACLE_BASE}"
echo "Oracle home is ${ORACLE_HOME}"
echo "Oracle SID is ${ORACLE_SID}"
PATH=$ORACLE_HOME/bin:$PATH; export PATH
REMOVE_ENTRY=$(sed -e '/^EMPTY.*/d' /etc/oratab)
echo "${REMOVE_ENTRY}" > /etc/oratab
sqlplus / as sysdba <<EOM
shutdown immediate;
quit
EOM
rm -Rf ${ORACLE_BASE}/admin/EMPTY ${ORACLE_BASE}/cfgtoollogs/dbca/EMPTY ${ORACLE_BASE}/oradata/EMPTY