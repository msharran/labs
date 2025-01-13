#!/bin/bash

COMMIT=$(git log -1 --pretty=%B | tr -d '\n')
RELEASE_TYPE=""

if [[ $COMMIT == "BREAKING CHANGE:"* ]]
then
  RELEASE_TYPE=major
elif [[ $COMMIT == "feat:"* ]]
then
  RELEASE_TYPE=minor
elif [[ $COMMIT == "fix:"* ]]
then
  RELEASE_TYPE=patch
fi

echo -n $RELEASE_TYPE
