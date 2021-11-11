#!/bin/sh

# Check working directory is the root directory of the PDKGO repository
if [ ! -d "docs/md" ]; then
    echo 'Please run this script from the root directory of the PDKGO project.'
    exit 1
fi
# Check hugo extended is installed
if ! hugo version | grep -q 'extended'; then
    echo 'The "extended" version of hugo was not found, please install it.'
    exit
fi
# Check git is installed
if ! type "git" > /dev/null; then
    echo "Git version control was not found, please install it."
    exit
fi
# Check if npm is installed
if ! type "npm" > /dev/null; then
    echo "NPM was not found, please install it."
    exit
fi

if [ ! -d "docs/site" ]; then
    git clone https://github.com/puppetlabs/devx.git docs/site
    cd docs/site
    echo "replace github.com/puppetlabs/pdkgo/docs/md => ../md" >> go.mod
    npm install
else
    cd docs/site
    git pull
    hugo mod clean
fi

# Check for -D flag to see if user wants to run the site with draft pages displayed
flag=''
while getopts 'D' opt; do
    case $opt in
        D) flag='-D' ;;
    esac
done

git submodule update --init --recursive --depth 50
hugo mod get
hugo server $flag
