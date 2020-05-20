#!/bin/bash

curl -s -X PUT "http://localhost:9200/datasets" -d @index.json
