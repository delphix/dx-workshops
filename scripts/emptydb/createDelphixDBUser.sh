#!/bin/sh
#
# Copyright (c) 2019 by Delphix. All rights reserved.
#

#
# Echo error messages before exiting with error code 1
#
die()
{
	for error_message in "$@"; do
		echo "$error_message"
	done
	exit 1
}

#
# Get a specific value from SQL*Plus ouput.
# The format to look for is <ValueName>=<value>
#
# Parameters
# 1. The SQL*Plus ouptut
# 2. The ValueName
#
# If the value is found successfully, the value is echoed to the stdout
#
get_value_from_sqlplus_output()
{
	SQLPLUS_OUTPUT=$1
	VALUE_NAME=$2
	echo "$SQLPLUS_OUTPUT" | sed -n "s/.*$VALUE_NAME=\(.*\)/\1/p"
}

gen_kccfe_priv() {
	cat >>$DBUSER_CREATION_SQL <<-EOF
		create or replace view v_x\$kccfe as select * from x\$kccfe;
		grant select on v_x\$kccfe to $USERNAME;
		create synonym $USERNAME.x\$kccfe for v_x\$kccfe;
	EOF
}

gen_select_any_dictionary_priv() {
	echo "grant select any dictionary to $USERNAME;" >>$DBUSER_CREATION_SQL
	gen_kccfe_priv
}

gen_view_specific_priv() {
	cat >>$DBUSER_CREATION_SQL <<-EOF
		grant select on v_\$active_instances to $USERNAME;
		grant select on v_\$archive_dest to $USERNAME;
		grant select on v_\$archive_dest_status to $USERNAME;
		grant select on v_\$archived_log to $USERNAME;
		grant select on v_\$backup to $USERNAME;
		grant select on v_\$backup_datafile to $USERNAME;
		grant select on v_\$backup_piece to $USERNAME;
		grant select on v_\$backup_set to $USERNAME;
		grant select on v_\$controlfile to $USERNAME;
		grant select on v_\$controlfile_record_section to $USERNAME;
		grant select on v_\$database to $USERNAME;
		grant select on v_\$database_incarnation to $USERNAME;
		grant select on v_\$datafile to $USERNAME;
		grant select on v_\$datafile_header to $USERNAME;
		grant select on v_\$instance to $USERNAME;
		grant select on v_\$license to $USERNAME;
		grant select on v_\$log to $USERNAME;
		grant select on v_\$logfile to $USERNAME;
		grant select on v_\$option to $USERNAME;
		grant select on v_\$parameter to $USERNAME;
		grant select on v_\$parameter2 to $USERNAME;
		grant select on gv_\$rman_configuration to $USERNAME;
		grant select on v_\$rman_configuration to $USERNAME;
		grant select on gv_\$session to $USERNAME;
		grant select on v_\$sqltext to $USERNAME;
		grant select on v_\$standby_log to $USERNAME;
		grant select on v_\$sysstat to $USERNAME;
		grant select on v_\$tablespace to $USERNAME;
		grant select on dba_tablespaces to $USERNAME;
		grant select on dba_temp_files to $USERNAME;
		grant select on v_\$tempfile to $USERNAME;
		grant select on v_\$thread to $USERNAME;
		grant select on gv_\$thread to $USERNAME;
		grant select on v_\$transaction to $USERNAME;
		grant select on v_\$version to $USERNAME;
		grant select on v_\$nls_parameters to $USERNAME;
		create synonym $USERNAME.v\$active_instances for v_\$active_instances;
		create synonym $USERNAME.v\$archive_dest for v_\$archive_dest;
		create synonym $USERNAME.v\$archive_dest_status for v_\$archive_dest_status;
		create synonym $USERNAME.v\$archived_log for v_\$archived_log;
		create synonym $USERNAME.v\$backup for v_\$backup;
		create synonym $USERNAME.v\$backup_datafile for v_\$backup_datafile;
		create synonym $USERNAME.v\$backup_piece for v_\$backup_piece;
		create synonym $USERNAME.v\$backup_set for v_\$backup_set;
		create synonym $USERNAME.v\$controlfile for v_\$controlfile;
		create synonym $USERNAME.v\$controlfile_record_section for v_\$controlfile_record_section;
		create synonym $USERNAME.v\$database for v_\$database;
		create synonym $USERNAME.v\$database_incarnation for v_\$database_incarnation;
		create synonym $USERNAME.v\$datafile for v_\$datafile;
		create synonym $USERNAME.v\$datafile_header for v_\$datafile_header;
		create synonym $USERNAME.v\$instance for v_\$instance;
		create synonym $USERNAME.v\$license for v_\$license;
		create synonym $USERNAME.v\$log for v_\$log;
		create synonym $USERNAME.v\$logfile for v_\$logfile;
		create synonym $USERNAME.v\$option for v_\$option;
		create synonym $USERNAME.v\$parameter for v_\$parameter;
		create synonym $USERNAME.v\$parameter2 for v_\$parameter2;
		create synonym $USERNAME.gv\$rman_configuration for gv_\$rman_configuration;
		create synonym $USERNAME.v\$rman_configuration for v_\$rman_configuration;
		create synonym $USERNAME.gv\$session for gv_\$session;
		create synonym $USERNAME.v\$sqltext for v_\$sqltext;
		create synonym $USERNAME.v\$standby_log for v_\$standby_log;
		create synonym $USERNAME.v\$sysstat for v_\$sysstat;
		create synonym $USERNAME.v\$tablespace for v_\$tablespace;
		create synonym $USERNAME.v\$tempfile for v_\$tempfile;
		create synonym $USERNAME.v\$thread for v_\$thread;
		create synonym $USERNAME.gv\$thread for gv_\$thread;
		create synonym $USERNAME.v\$transaction for v_\$transaction;
		create synonym $USERNAME.v\$version for v_\$version;
		create synonym $USERNAME.v\$nls_parameters for v_\$nls_parameters;

		-- check for v\$active_session_history
		declare
			exist int;
		begin
			select count(*) into exist from v\$fixed_table where name='V\$ACTIVE_SESSION_HISTORY';
			if exist = 1 then
				execute immediate 'grant select on v_\$active_session_history to $USERNAME';
				execute immediate 'create synonym $USERNAME.v\$active_session_history for v_\$active_session_history';
			end if;
		end;
		/

		-- check for v\$block_change_tracking
		declare
			exist int;
		begin
			select count(*) into exist from v\$fixed_table where name='V\$BLOCK_CHANGE_TRACKING';
			if exist = 1 then
				execute immediate 'grant select on v_\$block_change_tracking to $USERNAME';
				execute immediate 'create synonym $USERNAME.v\$block_change_tracking for v_\$block_change_tracking';
			end if;
		end;
		/

		-- check for v\$dnfs_channels
		declare
			exist int;
		begin
			select count(*) into exist from v\$fixed_table where name='V\$DNFS_CHANNELS';
			if exist = 1 then
				execute immediate 'grant select on v_\$dnfs_channels to $USERNAME';
				execute immediate 'create synonym $USERNAME.v\$dnfs_channels for v_\$dnfs_channels';
			end if;
		end;
		/

		-- check for v\$system_fix_control
		declare
			exist int;
		begin
			select count(*) into exist from v\$fixed_table where name='V\$SYSTEM_FIX_CONTROL';
			if exist = 1 then
				execute immediate 'grant select on v_\$system_fix_control to $USERNAME';
				execute immediate 'create synonym $USERNAME.v\$system_fix_control for v_\$system_fix_control';
			end if;
		end;
		/

		-- check for v\$transportable_platform
		declare
			exist int;
		begin
			select count(*) into exist from v\$fixed_table where name='V\$TRANSPORTABLE_PLATFORM';
			if exist = 1 then
				execute immediate 'grant select on v_\$transportable_platform to $USERNAME';
				execute immediate 'create synonym $USERNAME.v\$transportable_platform for v_\$transportable_platform';
			end if;
		end;
		/

		-- check for cdb_tablespaces
		declare
			exist int;
		begin
			select count(*) into exist from dba_catalog where table_name='CDB_TABLESPACES';
			if exist > 0 then
				execute immediate 'grant select on cdb_tablespaces to $USERNAME';
			end if;
		end;
		/

		-- check for v\$containers
		declare
			exist int;
		begin
			select count(*) into exist from v\$fixed_table where name='V\$CONTAINERS';
			if exist = 1 then
				execute immediate 'grant select on v_\$containers to $USERNAME';
				execute immediate 'grant select on gv_\$containers to $USERNAME';
				execute immediate 'create synonym $USERNAME.v\$containers for v_\$containers';
				execute immediate 'create synonym $USERNAME.gv\$containers for gv_\$containers';
			end if;
		end;
		/

		-- check for v\$pdb_incarnation
		declare
			exist int;
		begin
			select count(*) into exist from v\$fixed_table where name='V\$PDB_INCARNATION';
			if exist = 1 then
				execute immediate 'grant select on v_\$pdb_incarnation to $USERNAME';
				execute immediate 'create synonym $USERNAME.v\$pdb_incarnation for v_\$pdb_incarnation';
			end if;
		end;
		/

		EOF
	gen_kccfe_priv
}

echo
echo "This script will create the Delphix database user"
echo

if [ -z "$ORACLE_SID" ]; then
	echo "Please enter the Oracle SID of the database instance where the DB user is to be created:"
else
	echo "Database User will be created in instance <$ORACLE_SID>"
	echo "If you would like to create the DB user in a different Database Instance,"
	echo "please enter the desired SID or press <Enter> to proceed."
fi
read SID

if [ -n "$SID" ]; then
	ORACLE_SID=$SID; export ORACLE_SID
fi

#
# If the OS user doesn't have sysdba privilege, stop now.
#
sqlplus "/ as sysdba" </dev/null 2>&1 | grep "Connected" >/dev/null ||
	die "Please run this script with an OS user that has SYSDBA privilege for database instance <$ORACLE_SID>."

#
# Find out the database version
#
sqlplus_output=`sqlplus -S -R 3 "/ as sysdba" <<-EOF
	set heading off pagesize 0 echo off define off feedback off verify off trimout on;
	select 'Version=' || version from v\\$instance;
	exit;
EOF`
DATABASE_VERSION=`get_value_from_sqlplus_output "$sqlplus_output" Version`
DATABASE_MAJOR_VERSION=`echo $DATABASE_VERSION | awk -F"." '{print $1}'`

#
# Test whether we got a legitimate number for the database version
#
expr "$DATABASE_MAJOR_VERSION" + 1 >/dev/null 2>&1 ||
	die "$sqlplus_output" "Unable to get database version from database instance <$ORACLE_SID>."


#
# For Oracle 12c and later version we need to handle container databases and pluggable databases
#
IS_CDB=false
if [ "$DATABASE_MAJOR_VERSION" -gt 11 ]; then
	sqlplus_output=`sqlplus -S -R 3 "/ as sysdba" <<-EOF
		set heading off pagesize 0 echo off define off feedback off verify off trimout on;
		select 'CDB=' || CDB from v\\$database;
		exit;
	EOF`
	[ `get_value_from_sqlplus_output "$sqlplus_output" CDB` = "YES" ] && IS_CDB=true
fi

PDB_NAME=
if $IS_CDB; then
	echo "<$ORACLE_SID> is a container database instance."
	echo "Pleaser enter the pluggable database name for which to create the Delphix database user,"
	echo "or press <Enter> to create the Delphix database user for CDB\$ROOT:"
	read PDB_NAME
fi

echo "Please enter the Delphix database user name to be created:"
$IS_CDB && [ -z "$PDB_NAME" ] && echo "(CDB\$ROOT user name must start with C## or c##.)"
read USERNAME

#
# Verify that CDB$ROOT user name must start with C## or c##
#
if $IS_CDB; then
	#
	# Check whether the user name start with c## or not
	#
	case $USERNAME in
		[cC]##*) COMMON_USER_PREFIX=true ;;
		*) COMMON_USER_PREFIX=false ;;
	esac

	[ -n "$PDB_NAME" ] && $COMMON_USER_PREFIX &&
		die "Error: Database user for pluggable databases cannot start with C## or c##."

	[ -z "$PDB_NAME" ] && [ "$COMMON_USER_PREFIX" = "false" ] &&
		die "Error: Database user for CDB\$ROOT must start with C## or c##."
fi

echo "Please enter the password for database user \"$USERNAME\":"
echo "(Oracle does not support passwords starting with a number.)"
stty -echo
read PASSWORD
stty echo

echo
echo -n "Grant SELECT ANY DICTIONARY privilege to \"$USERNAME\"? [(y)/n]: "
read response
[ -z "$response" ] && response="y"

case "$response" in
	[yY]|[yY][eE][sS]) anyselect=1;;
	* ) anyselect=0;;
esac

#
# Generate a SQL script to create the Delphix DB user with necessary privileges.
# If the user does not want 'select any dictionary' system privilege handed out.
# then grant select access to required v$ view base tables and create the synonyms.
#
echo
DBUSER_CREATION_SQL=createDelphixDBUser.sql
echo "Generating SQL script: $DBUSER_CREATION_SQL"

>$DBUSER_CREATION_SQL

[ -n "$PDB_NAME" ] && echo "alter session set container=$PDB_NAME;" >>$DBUSER_CREATION_SQL

echo "create user $USERNAME identified by $PASSWORD;" >>$DBUSER_CREATION_SQL
echo "grant create session to $USERNAME;" >>$DBUSER_CREATION_SQL

$IS_CDB && [ -z "$PDB_NAME" ] &&
	echo "alter user $USERNAME set container_data=all container=current;" >>$DBUSER_CREATION_SQL

if [ $anyselect -eq 1 ]; then
	gen_select_any_dictionary_priv
else
	gen_view_specific_priv
fi

echo exit >>$DBUSER_CREATION_SQL

echo
echo "Executing SQL script: $DBUSER_CREATION_SQL"
sqlplus -s "/ as sysdba" "@$DBUSER_CREATION_SQL"
