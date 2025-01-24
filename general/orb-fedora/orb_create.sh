#!/bin/bash

set -e

orb delete fedora
orb create -a amd64 -p -u msharran fedora
orb -m fedora push install.sh /home/msharran/install.sh
