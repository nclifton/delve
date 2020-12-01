#!/bin/bash

SCRIPTDIR=$(/usr/bin/dirname $0)
if [ "$SCRIPTDIR" != "." ]; then
  cd $SCRIPTDIR
fi
SCRIPTDIR=$(/bin/pwd)
ROOTDIR="$(cd ../../ >/dev/null 2>&1 && pwd)"

# Load in service paths based on whether a directory contains an infra folder or not
cd ${ROOTDIR}
service_paths=$(find ./ -type d -name infra | sed -e 's/\/infra$//' | sed -e 's/\.\///')
cd ${SCRIPTDIR}

# Update all Helm Dependencies (this takes a while)
for service_path in ${service_paths}; do
    infra="${ROOTDIR}/${service_path}/infra"
    dir=$(find ${infra} -mindepth 1 -maxdepth 1 -type d | egrep -v "values")

    helm dependency update ${dir} --skip-refresh
done
