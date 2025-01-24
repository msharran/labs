#!/bin/bash

orbctl -m fedora push installer.sh /home/msharran/installer.sh
orbctl -m fedora run bash /home/msharran/installer.sh
