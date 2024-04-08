#!/bin/bash

set -e

TEST_DIR=/tmp/bramble-cli-test/
GIT_TEST_DIR=$TEST_DIR/git-test

BIN=$PWD/bin/bramble

mkdir -p $GIT_TEST_DIR

testCommandSad() {
   if "$@"; then 
      false
   else
      true
   fi
}

testCommandHappy() {
    "$@" 
}

cleanUpOnFailure() {
   "$@" || (rm -rf $TEST_DIR && exit 1)
}

printTestMSG() {
   printf "\n\n-----test %s-----\n\n" "$@"
}

pushd $GIT_TEST_DIR || exit

   printTestMSG "no directory specified"
   cleanUpOnFailure testCommandSad "$BIN" init

   printTestMSG ". empty path -- SAD"
   cleanUpOnFailure testCommandSad "$BIN" init .

   printTestMSG "invalid git format -- SAD"
   touch .git
   cleanUpOnFailure testCommandSad "$BIN" init .
   rm -f .git

   printTestMSG ". path, git -- HAPPY"
   git init . --quiet
   testCommandHappy "$BIN" init .

   printf "\n\nAll tests passing\n\n"

popd || exit

rm -rf $TEST_DIR
