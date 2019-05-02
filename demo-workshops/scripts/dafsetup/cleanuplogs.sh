#!/bin/bash
#
# Copyright (c) 2019 by Delphix. All rights reserved.
#
PGSQLBIN=/usr/pgsql-9.6/bin
ARCHIVELOGDIR=/var/lib/pgsql/9.6/backups

export PATH=${PGSQLBIN}:${PATH}

file="$(find ${ARCHIVELOGDIR}/ -type f -printf "%C@ %f\n" | sort -n | tail -n 1 | awk '{print $NF}')"
pg_archivecleanup ${ARCHIVELOGDIR}/ $file