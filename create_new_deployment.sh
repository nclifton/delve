#!/bin/bash
SCRIPTDIR=`/usr/bin/dirname $0`
if [ "$SCRIPTDIR" != "." ]; then
  cd $SCRIPTDIR
fi

command -v git >/dev/null 2>&1 || { echo >&2 "Command 'git' is not installed or is not executable, aborting."; exit 1; }

# 0. Any uncommitted changes?
if [[ ! -z $(git status --porcelain) ]]; then
    echo "You have uncommitted changes, commit before running this script"
    exit 1
fi

echo "1. Fetch origin"
git fetch origin

echo "2. Checkout develop branch"
git checkout develop
git pull

echo "3. Delete current local main branch"
git branch -D main

echo "4. Create new main branch"
git checkout -b main

echo "5. Force push to main branch"
git push -f origin main

echo "6. Checkout develop branch"
git checkout develop

echo "Deployment successfully pushed"
echo "Once built, the deployment will automatically be deployed to QA"
