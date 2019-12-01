#!/usr/bin/env bash
set -e

STARTTIME=$(date +%s)
NOW=$(date +"%m-%d-%Y %T")

function ENDTIME {
        ENDTIME=$(date +%s)
        echo "It took $(($ENDTIME - $STARTTIME)) seconds to complete ${SCRIPT}"
}

function ERROR {
        ENDTIME
        exit 1
}

{
sed -i 's|host\ tooling|host\ \*|' ~/.ssh/config || ERROR
./add_patient.sh || ERROR
./execute_data_pipeline.sh || ERROR
./provision_vdb.sh || ERROR
./push_notes_field.sh || ERROR
./fix_notes.sh || ERROR
./merge_notes.sh || ERROR
./fix_test.sh || ERROR
./merge_testfix.sh || ERROR
./push_to_prod.sh || ERROR
} 2>&1 | tee "$(date '+%F-%H%M%S').log"

exit ${PIPESTATUS[0]}