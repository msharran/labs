.PHONY: build rebuild

build:
	./build_vm.sh

rebuild:
	orbctl -m fedora push installer.sh /home/msharran/installer.sh
	orbctl -m fedora run /home/msharran/installer.sh
	infocmp -x | ssh msharran@fedora@orb -- tic -x -
