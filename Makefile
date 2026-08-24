.PHONY: install update submodules

install:
	@which python3 > /dev/null 2>&1 || (echo "Error: python3 is not installed" && exit 1)
	@test -d dotfiles-venv || $(MAKE) create-venv
	@dotfiles-venv/bin/pip show ansible > /dev/null 2>&1 || dotfiles-venv/bin/pip install ansible
	@sudo -v
	dotfiles-venv/bin/ansible-playbook ansible/main.yml -u $(USER) --ask-become-pass

update:
	git pull
	git submodule update --remote --recursive

submodules:
	git submodule update --init --recursive

create-venv:
	python3 -m venv dotfiles-venv

a-venv:
	@echo "Run: source dotfiles-venv/bin/activate"
