## RUN
git clone --recurse-submodules git@github.com:shayne651/dotfiles.git

ansible-playbook ansible/main.yml -u $USER --ask-become-pass


## DEPENDENCIES
ansible
git
