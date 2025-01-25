#!/usr/bin/env bash

run_cmd() {
    eval $@
    if [[ $? != 0 ]]; then 
        exit 1
    fi
}

if orbctl list 2>/dev/null | grep -q fedora
then
    echo "* Found fedora VM, deleting it"
    run_cmd orbctl delete fedora 
fi

echo "* Creating fedora VM"
run_cmd orbctl create -a amd64 -p -u msharran fedora

echo "* Pushing Installer into the VM"
run_cmd orbctl -m fedora push installer.sh /home/msharran/installer.sh

echo "* Running Installer Script"
run_cmd orbctl -m fedora run bash /home/msharran/installer.sh

echo "* Copying term info"
run_cmd "infocmp -x | ssh msharran@fedora@orb -- tic -x -"
