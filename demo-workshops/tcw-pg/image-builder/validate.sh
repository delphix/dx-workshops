#!/usr/bin/env bash
#
# Copyright (c) 2019 by Delphix. All rights reserved.
#
source $(dirname "${BASH_SOURCE[0]}")/library.sh
trap "cleanup" SIGINT

S3_BINARIES=(Plugin_postgreSQL_1.3.0.zip datical_admin.lic DaticalDB-linux.gtk.x86_64-5.2.5347.jar linuxx64_12201_database.zip oracle-instantclient12.2-basic-12.2.0.1.0-1.x86_64.rpm oracle-instantclient12.2-tools-12.2.0.1.0-1.x86_64.rpm oracle-instantclient12.2-jdbc-12.2.0.1.0-1.x86_64.rpm oracle-instantclient12.2-sqlplus-12.2.0.1.0-1.x86_64.rpm)

function HELP() {
    echo -e "\nPlease visit https://github.com/delphix/dx-workshops/blob/master/demo-workshops/tcw/docs/building.md to validate the pre-requisites."
    exit 1
}

echo "Validating environment variables"

: "${AWS_ACCESS_KEY_ID:?Unset. This is a mandatory value}"
: "${AWS_KEYNAME:?Unset. This is a mandatory value}"
: "${AWS_REGION:?Unset. This is a mandatory value}"
: "${AWS_SECRET_ACCESS_KEY:?Unset. This is a mandatory value}"
: "${AWS_SUBNET_ID:?Unset. This is a mandatory value}"
: "${AWS_VPC_ID:?Unset. This is a mandatory value}"
: "${DELPHIX_ADMIN_PASSWORD:?Unset. This is a mandatory value}"
: "${GUACADMIN_PASSWORD:?Unset. This is a mandatory value}"
: "${GUACAMOLE_DB_PASSWORD:?Unset. This is a mandatory value}"
: "${GUAC_USER_PASSWORD:?Unset. This is a mandatory value}"
: "${MARIADB_ROOT_PASSWORD:?Unset. This is a mandatory value}"
: "${S3_BUCKET:?Unset. This is a mandatory value}"
: "${S3_OBJECT_PATH:?Unset. This is a mandatory value}"
: "${VNC_DEFAULT_PASSWORD:?Unset. This is a mandatory value}"
: "${DELPHIX_VERSION:?Unset. This is a mandatory value}"

#Check that the Delphix AMI is shared with the Account and Region
echo -e "\nChecking that the Delphix Engine ${DELPHIX_VERSION} is shared with this account in the ${AWS_REGION} region..."
if [[ -n "$(aws ec2 --region ${AWS_REGION} describe-images --filters "Name=name,Values=Delphix Engine ${DELPHIX_VERSION}" --owner "180093685553" --output text)" ]]; then
    echo "Found Delphix Engine ${DELPHIX_VERSION}"
else
    echo "Could not find Delphix Engine ${DELPHIX_VERSION} shared with this region"
    echo "For registered customers, you can login to https://download.delphix.com and share the ami with your account and region"
    echo -e "For everyone else, please contact your account representative to speak about access.\n"
    FAIL=1
fi

#Check that the s3 bucket with files are present
echo -e "\nChecking access to the S3 bucket and the binaries are present..."
S3_LIST=$(aws s3 --region ${AWS_REGION} ls s3://${S3_BUCKET}${S3_OBJECT_PATH}/|awk '{print $4}')

if [[ -z $S3_LIST ]]; then 
    echo -e "\nDid not find the S3 bucket path specified was invalid/empty: ${S3_BUCKET}${S3_OBJECT_PATH}\n"
    FAIL=1
else
    for each in ${S3_BINARIES[@]}; do
        if [[ $S3_LIST == *"${each}"* ]]; then
            echo "Found $each"
        else
            echo "Did not find $each"
            FAIL=1
        fi
    done
fi

echo $FAIL

[[ ${FAIL} ]] && HELP

