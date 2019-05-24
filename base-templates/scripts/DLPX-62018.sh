#
# Copyright (c) 2019 by Delphix. All rights reserved.
#
#https://jira.delphix.com/browse/DLPX-62018
postgresql_data_dir=$1
sudo sed -i -e "s|\(port = \)\(\.*'\)\(.*\)\(.*'\)|\1\3|" ${postgresql_data_dir}/postgresql.conf