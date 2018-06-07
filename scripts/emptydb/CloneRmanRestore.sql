SET VERIFY OFF
connect "SYS"/"&&sysPassword" as SYSDBA
set echo on
spool /home/oracle/emptydb/CloneRmanRestore.log append
startup nomount pfile="/home/oracle/emptydb/init.ora";
@/home/oracle/emptydb/rmanRestoreDatafiles.sql;
spool off
