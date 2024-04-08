#!/bin/bash

set -e

TEST_DIR=/tmp/bramble-cli-test/
GIT_TEST_DIR=$TEST_DIR/git-test
EMPTY_TEST_DIR=$TEST_DIR/empty-test

BIN=$PWD/bin/bramble

mkdir -p $GIT_TEST_DIR
mkdir -p $EMPTY_TEST_DIR

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

pushd $GIT_TEST_DIR
   printf "\n\n-----test . path, git -- sad-----\n\n"
   testCommandSad "$BIN" init .

   printf "\n\n-----test invalid git format -- sad-----\n\n"
   touch .git 
   testCommandSad "$BIN" init .
   rm .git

   printf "\n\n------test . path, git -- happy-----\n\n"
   git init .

   testCommandHappy "$BIN" init .
popd 

rm -rf $TEST_DIR
