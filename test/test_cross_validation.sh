#!/usr/bin/env bash
set -e

echo '=== RUN   Test Cross Validation'

# Configuration
EPSILON=0.005
LOCATION=$(dirname "$0")

# Build executable
go build -o ${LOCATION}/gorse ${LOCATION}/../cmd/gorse.go

# Test executable
declare RESULT=($(${LOCATION}/gorse cv svd --load-builtin ml-100k \
    --eval-rmse \
    --eval-mae \
    --set-n-epochs 100 \
    --set-reg 0.1 \
    --set-lr 0.01 \
    --set-n-factors 50 \
    --set-init-mean 0 \
    --set-init-std 0.001 | grep -P '\d+\.\d+(?=\()' -o))

echo $(date +'%Y/%m/%d %H:%M:%S') RMSE = ${RESULT[0]}, MAE = ${RESULT[1]}

# Check RMSE
if [[ $(echo "${RESULT[0]}-0.90728<${EPSILON}" | bc) != 1 ]]; then
    echo "--- FAIL  unexpected RMSE: ${RESULT[0]} - 0.90728 >= ${EPSILON}"
    exit 1
fi

# Check MAE
if [[ $(echo "${RESULT[1]}-0.71508<${EPSILON}" | bc) != 1 ]]; then
    echo "--- FAIL  unexpected MAE: ${RESULT[1]} - 0.71508 >= ${EPSILON}"
    exit 1
fi

# Clear
bash ${LOCATION}/clean_file.sh

# Print success
echo '--- PASS  Test Cross Validation'
