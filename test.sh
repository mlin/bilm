#!/bin/bash -e

SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ] ; do SOURCE="$(readlink "$SOURCE")"; done
srcdir="$( cd -P "$( dirname "$SOURCE" )" && pwd )" 

testdir=$(mktemp --tmpdir -d bilm_test.XXXXXX)

cd $srcdir/bilm_add
go build
cp bilm_add $testdir

cd $srcdir/bilm_query
go build
cp bilm_query $testdir

cd $testdir
echo "The quick brown fox jumps over the lazy dog" | ./bilm_add test.bilm
echo "Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.
Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.
Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur.
Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum." | ./bilm_add test.bilm

if  ./bilm_query test.bilm DEADBEEF ; then
	echo "should have failed: $CMD"
	exit 1
fi

if [ $(./bilm_query test.bilm "THE quicK brown fox" | wc -l) -ne "1" ]; then
	echo 'shouild have produced one result: ./bilm_query test.bilm "THE quicK brown fox"'
	exit 1
fi

if [ $(./bilm_query test.bilm "jumps over the * dog" | wc -l) -ne "1" ]; then
	echo 'shouild have produced one result: ./bilm_query test.bilm "jumps over the * dog"'
	exit 1
fi

if [ $(./bilm_query test.bilm "in" | wc -l) -ne "2" ]; then
	echo 'shouild have produced two results: ./bilm_query test.bilm "in"'
	exit 1
fi

rm -rf $testdir
echo "PASS"
