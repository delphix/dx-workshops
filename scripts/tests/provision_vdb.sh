#!/bin/bash
echo "###Validate MASKED VDB###"
set -e
DEIP="virtualizationengine"
USER="delphix_admin"
PASS="Landshark-12"

function JOB_WAIT {
    JOBID=${1}
    until [[ "$(curl -s "http://${DEIP}/resources/json/delphix/job/${JOBID}" \
        -b ~/cookies.txt -H "Content-Type: application/json" \
        | jq -r .result.jobState)" != "RUNNING" ]]; do 
        echo "Waiting for ${JOBID}"
        sleep 3
    done

    JOB_STATE=$(curl -s "http://${DEIP}/resources/json/delphix/job/${JOBID}" \
        -b ~/cookies.txt -H "Content-Type: application/json" \
        | jq -r .result.jobState)

    echo "${JOBID} ${JOB_STATE}"
    if [[ "${JOB_STATE}" != "COMPLETED" ]]; then
        exit 1
    fi
}

# 1) Create Delphix API Session
echo -e "\n\nCreating Session\n"
curl -sX POST -k --data @- http://${DEIP}/resources/json/delphix/session \
    -c ~/cookies.txt -H "Content-Type: application/json" <<EOF
{
    "type": "APISession",
    "version": {
        "type": "APIVersion",
        "major": 1,
        "minor": 9,
        "micro": 0
    }
}
EOF


# 2) Delphix Login
echo -e "\n\nLogging in\n"
curl -sX POST -k --data @- http://${DEIP}/resources/json/delphix/login \
    -b ~/cookies.txt -c ~/cookies.txt -H "Content-Type: application/json" <<EOF
{
    "type": "LoginRequest",
    "username": "${USER}",
    "password": "${PASS}"
}
EOF

# echo "Refreshing Environments"
# ENVJOBID1=$(curl -sX POST \
#         "http://${DEIP}/resources/json/delphix/environment/UNIX_HOST_ENVIRONMENT-1/refresh" \
#         -b ~/cookies.txt -H 'Content-Type: application/json'| jq -r .job)
# ENVJOBID2=$(curl -sX POST \
#         "http://${DEIP}/resources/json/delphix/environment/UNIX_HOST_ENVIRONMENT-2/refresh" \
#         -b ~/cookies.txt -H 'Content-Type: application/json'| jq -r .job)

# JOB_WAIT $ENVJOBID1
# JOB_WAIT $ENVJOBID2

ENVUSER=$(curl -s "http://${DEIP}/resources/json/delphix/environment" \
    -b ~/cookies.txt -H 'Content-Type: application/json' | \
    jq -r ".result[]| select(.name==\"proddb\").primaryUser")

echo $ENVUSER

ENVREF=$(curl -s "http://${DEIP}/resources/json/delphix/environment" \
    -b ~/cookies.txt -H 'Content-Type: application/json' | \
    jq -r ".result[]| select(.name==\"proddb\").reference")

echo $ENVREF

if [[ $(ssh centos@proddb sudo -i -u oracle whoami) ]]; then
    echo -e "\nOracle it is\n"
    
    REPOREF=$(curl -s "http://${DEIP}/resources/json/delphix/sourceconfig?environment=${ENVREF}" \
    -b ~/cookies.txt -H 'Content-Type: application/json' | \
    jq -r ".result[]| select(.name==\"PATMM\").repository")

    echo ${REPOREF}

    CDBCONFIG=$(curl -s "http://${DEIP}/resources/json/delphix/sourceconfig?environment=${ENVREF}" \
    -b ~/cookies.txt -H 'Content-Type: application/json' | \
    jq -r ".result[]| select(.name==\"PATMM\").cdbConfig")

    echo ${CDBCONFIG}

    VDBJOB=$(curl -sX POST \
        "http://${DEIP}/resources/json/delphix/database/provision" \
        -b ~/cookies.txt -H 'Content-Type: application/json' \
        -d @- <<-EOF
        {
            "sourceConfig": {
                "databaseName": "test",
                "cdbConfig": "${CDBCONFIG}",
                "nonSysUser": null,
                "nonSysCredentials": null,
                "linkingEnabled": true,
                "environmentUser": "${ENVUSER}",
                "repository": "${REPOREF}",
                "type": "OraclePDBConfig"
            },
            "source": {
                "operations": {
                    "configureClone": [

                    ],
                    "preRefresh": [

                    ],
                    "postRefresh": [

                    ],
                    "preRollback": [

                    ],
                    "postRollback": [

                    ],
                    "preSnapshot": [

                    ],
                    "postSnapshot": [

                    ],
                    "preStart": [

                    ],
                    "postStart": [

                    ],
                    "preStop": [

                    ],
                    "postStop": [

                    ],
                    "type": "VirtualSourceOperations"
                },
                "customEnvVars": [

                ],
                "allowAutoVDBRestartOnHostReboot": false,
                "mountBase": "/mnt/provision",
                "logCollectionEnabled": false,
                "name": "test",
                "type": "OracleVirtualPdbSource"
            },
            "container": {
                "diagnoseNoLoggingFaults": true,
                "preProvisioningEnabled": false,
                "sourcingPolicy": {
                    "logsyncMode": "UNDEFINED",
                    "logsyncInterval": 5,
                    "logsyncEnabled": false,
                    "type": "OracleSourcingPolicy"
                },
                "performanceMode": "DISABLED",
                "group": "GROUP-4",
                "name": "test",
                "type": "OracleDatabaseContainer"
            },
            "timeflowPointParameters": {
				"type": "TimeflowPointSemantic",
				"location": "LATEST_SNAPSHOT",
				"container": "ORACLE_DB_CONTAINER-4"
			},
            "masked": false,
            "type": "OracleMultitenantProvisionParameters"
        }
EOF
)
elif [[ $(ssh centos@proddb sudo -i -u postgres whoami) ]]; then
    echo -e "\nPostgres it is\n"
    REPOREF=$(curl -s "http://${DEIP}/resources/json/delphix/sourceconfig?environment=${ENVREF}" \
    -b ~/cookies.txt -H 'Content-Type: application/json' | \
    jq -r ".result[]| select(.path==\"/mnt/provision/patmm\").repository")
    
    echo ${REPOREF}
    
    VDBJOB=$(curl -sX POST \
        "http://${DEIP}/resources/json/delphix/database/provision" \
        -b ~/cookies.txt -H 'Content-Type: application/json' \
        -d @- <<-EOF
         {
            "container": {
                "sourcingPolicy": {
                    "logsyncEnabled": false,
                    "type": "SourcingPolicy"
                },
                "group": "GROUP-4",
                "name": "test",
                "type": "AppDataContainer"
            },
            "source": {
                "operations": {
                    "configureClone": [

                    ],
                    "preRefresh": [

                    ],
                    "postRefresh": [

                    ],
                    "preRollback": [

                    ],
                    "postRollback": [

                    ],
                    "preSnapshot": [

                    ],
                    "postSnapshot": [

                    ],
                    "preStart": [

                    ],
                    "postStart": [

                    ],
                    "preStop": [

                    ],
                    "postStop": [

                    ],
                    "type": "VirtualSourceOperations"
                },
                "parameters": {
                    "postgresPort": 5477,
                    "configSettingsStg": [
                        {
                            "propertyName": "listen_addresses",
                            "value": "*"
                        }
                    ]
                },
                "additionalMountPoints": [

                ],
                "allowAutoVDBRestartOnHostReboot": false,
                "logCollectionEnabled": false,
                "name": "test",
                "type": "AppDataVirtualSource"
            },
            "sourceConfig": {
                "path": "/mnt/provision/test",
                "name": "test",
                "repository": "${REPOREF}",
                "linkingEnabled": true,
                "environmentUser": "${ENVUSER}",
                "type": "AppDataDirectSourceConfig"
            },
            "timeflowPointParameters": {
				"type": "TimeflowPointSemantic",
				"location": "LATEST_SNAPSHOT",
				"container": "APPDATA_CONTAINER-2"
			},
            "masked": false,
            "type": "AppDataProvisionParameters"
        }
EOF
)
else
    echo "Nether expected user found. Exiting"
    exit 1
fi

echo ""

VDBJOBID=$(echo $VDBJOB | jq -r .job)

if [[ "$VDBJOBID" == "null" ]]; then
    echo "VDB wasn't created. Failing"
    echo $VDBJOB
    exit 1
fi

JOB_WAIT $VDBJOBID

echo -e "\nValidating Data\n" 
if [[ $(ssh centos@proddb sudo -i -u oracle whoami) ]]; then
    echo -e "\nOracle it is\n"
    OUTPUT=$(ssh centos@proddb sudo -i -u oracle sqlplus -s delphixdb/delphixdb@localhost:1521/test <<-EOF
    set pagesize 0 feedback off verify off heading off echo off;
    select lastname,city
    from patients
    where id = ( select max(id) from patients );
    quit;
EOF
    )

    if [[ "$(echo $OUTPUT| awk '{print $1}')" == "Bowen" || "$(echo $OUTPUT| awk '{print $2}')" != "Funkytown" ]]; then
        echo "Unexpected data"
        echo $OUTPUT
        exit 1
    fi
elif [[ $(ssh centos@proddb sudo -i -u postgres whoami) ]]; then
    echo -e "\Postgres it is\n"
    OUTPUT=$(psql -h proddb -p 5477 -U delphixdb dafdb -t -c 'select lastname,city  from patients where id = ( select max(id) from patients );')
    if [[ "$(echo $OUTPUT| awk '{print $1}')" == "Bowen" || "$(echo $OUTPUT| awk '{print $3}')" != "Funkytown" ]]; then
        echo "Unexpected data"
        echo $OUTPUT
        exit 1
    fi
fi


echo -e "\nData validated\n"

