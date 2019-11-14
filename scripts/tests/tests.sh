#!/usr/bin/env bash
set -e

{
 ./add_daf_user.sh
./execute_data_pipeline.sh
./push_notes_field.sh
./fix_notes.sh
./merge_feature.sh
./fix_test.sh
./merge_feature.sh
./push_to_prod.sh
} 2>&1 | tee "$(date '+%F-%H%M%S').log"