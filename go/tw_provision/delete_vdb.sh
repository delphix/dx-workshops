#/usr/local/env bash
ssh delphix_admin@delphix-tcw-virtualizationengine 'cd database; select "Patients Masked Master"; delete; set force=true;commit'
ssh delphix_admin@delphix-tcw-virtualizationengine 'cd database; select "Patients Prod"; delete; set force=true;commit'
ssh centos@delphix-tcw-source 'sudo pkill -f "/usr/pgsql-11/bin/postgres -D /var/lib/pgsql/staging/data"'
