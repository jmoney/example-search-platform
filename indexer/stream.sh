#!/bin/bash

for entry in `cat $1`; do
    curl -s -X PUT "http://localhost:8000/" -d $entry
    echo "Status Code: $?"
done