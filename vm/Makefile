SHELL = /bin/bash

# SSH "Host" added in ~/.ssh/config.
# This is already configured in my ssh config.
# But IP has to be configured when changing the VM
VM_SSH_HOST ?= vm

vm/bootstrap:
	$(MAKE) vm/secrets
	rsync -arP ./bootstrap.sh $(VM_SSH_HOST):~/bootstrap.sh && \
	ssh $(VM_SSH_HOST) "bash ~/bootstrap.sh"

vm/secrets:
	rsync -arP ~/.dotfiles-private $(VM_SSH_HOST):~/projects
	ssh $(VM_SSH_HOST) "cd ~/.dotfiles-private && make"
