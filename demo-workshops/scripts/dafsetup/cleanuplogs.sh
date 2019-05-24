#!/bin/bash
#
# Copyright (c) 2019 by Delphix. All rights reserved.
#
PGVER=REPLACE
PGSQLBIN=/usr/pgsql-${PGVER}/bin
ARCHIVELOGDIR=/var/lib/pgsql/${PGVER}/backups

export PATH=${PGSQLBIN}:${PATH}

file="$(find ${ARCHIVELOGDIR}/ -type f -printf "%C@ %f\n" | sort -n | tail -n 1 | awk '{print $NF}')"
pg_archivecleanup ${ARCHIVELOGDIR}/ $file