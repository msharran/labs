SHELL = /bin/bash

ORB_MACHINE ?= ubuntu
VM_USER ?= msharran
ANSIBLE_DIR := $(CURDIR)/ansible
ANSIBLE_INVENTORY := $(ORB_MACHINE),
ANSIBLE_PLAYBOOK := $(ANSIBLE_DIR)/ubuntu.yml
ANSIBLE_PLAYBOOK_CMD ?= ansible-playbook
ANSIBLE_EXTRA_VARS := ansible_host=orb ansible_user=ubuntu ansible_python_interpreter=/usr/bin/python3 target_user=$(VM_USER)

VAGRANT := vagrant
VAGRANT_DIR := $(CURDIR)/builder
HOSTOS := $(shell uname)
ARCH ?= $(shell uname -m | sed 's/x86_64/amd64/g; s/aarch64/arm64/g')

# SSH "Host" added in ~/.ssh/config.
# This is already configured in my ssh config.
# But IP has to be configured when changing the VM
VM_SSH_HOST ?= vm

# Commands prefixed with "host/" means they are executed on the host machine.
# Commands prefixed with "vm/" means they are executed on the VM.

host/bootstrap:
	$(MAKE) vm/secrets
	rsync -arP ./bootstrap.sh $(VM_SSH_HOST):~/bootstrap.sh && \
	ssh $(VM_SSH_HOST) "bash ~/bootstrap.sh"

host/secrets:
	rsync -arP ~/.dotfiles-private $(VM_SSH_HOST):~/projects
	ssh $(VM_SSH_HOST) "cd ~/.dotfiles-private && make"

.PHONY: ubuntu/install ubuntu/uninstall ubuntu/ansible-install ubuntu/ssh
.PHONY: vagrant-up vagrant-ssh vagrant-halt vagrant-destroy

ubuntu/install:
	@if orb list -q | grep -qx '$(ORB_MACHINE)'; then \
		echo "OrbStack machine '$(ORB_MACHINE)' already exists"; \
	else \
		orb create -a amd64 ubuntu $(ORB_MACHINE); \
	fi

ubuntu/uninstall:
	orb delete --force $(ORB_MACHINE)

ubuntu/ansible-install: ubuntu/install
	$(ANSIBLE_PLAYBOOK_CMD) -i $(ANSIBLE_INVENTORY) -e "$(ANSIBLE_EXTRA_VARS)" $(ANSIBLE_PLAYBOOK)

ubuntu/ssh:
	orb -m $(ORB_MACHINE)

vagrant-up: .vagrant-up
vagrant-ssh: .vagrant-ssh
vagrant-halt: .vagrant-halt
vagrant-destroy: .vagrant-destroy

.vagrant-%:
	VAGRANT_VAGRANTFILE=$(VAGRANT_DIR)/Vagrantfile-ubuntu \
		ARCH=$(ARCH) \
		HOSTOS=$(HOSTOS) \
		VM_USER=$(VM_USER) \
		$(VAGRANT) $*
