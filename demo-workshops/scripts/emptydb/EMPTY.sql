#
# Copyright (c) 2019 by Delphix. All rights reserved.
#
set verify off
ACCEPT sysPassword CHAR PROMPT 'Enter new password for SYS: ' HIDE
ACCEPT systemPassword CHAR PROMPT 'Enter new password for SYSTEM: ' HIDE
host /u01/app/oracle/product/11.2.0.4/ora_1/bin/orapwd file=/u01/app/oracle/product/11.2.0.4/ora_1/dbs/orapwEMPTY force=y
@/home/oracle/emptydb/CloneRmanRestore.sql
@/home/oracle/emptydb/cloneDBCreation.sql
@/home/oracle/emptydb/postScripts.sql
@/home/oracle/emptydb/lockAccount.sql
@/home/oracle/emptydb/postDBCreation.sql
SHUTDOWN IMMEDIATE;
STARTUP MOUNT;
ALTER DATABASE ARCHIVELOG;
ALTER DATABASE OPEN;
quit
