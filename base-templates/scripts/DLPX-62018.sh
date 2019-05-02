#
# Copyright (c) 2019 by Delphix. All rights reserved.
#
#https://jira.delphix.com/browse/DLPX-62018
sudo sed -i -e "s|\(port = \)\(\.*'\)\(.*\)\(.*'\)|\1\3|" /var/lib/pgsql/9.6/data/postgresql.conf