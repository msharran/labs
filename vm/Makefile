SHELL = /bin/bash

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
