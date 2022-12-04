#!/bin/bash

# A directory of test files to upload
TEST_FILES_DIR=test-upload-files
# All of the files downloadable by the server
TEST_UPLOAD_DIR=files

for f in $(ls $TEST_FILES_DIR); do
  curl -X POST --form file="@$TEST_FILES_DIR/$f" http://localhost:8080/upload
  echo ""
done

for f in $(ls $TEST_UPLOAD_DIR); do
  curl -X GET http://localhost:8080/download/$f
  echo ""
done