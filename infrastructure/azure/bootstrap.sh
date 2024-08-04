#!/bin/bash

# Clone your Ansible repository
git clone https://github.com/endalk200/termflow.git /etc/termflow

# Run your initial Ansible playbook
ansible-playbook /etc/termflow/infrastructure/azure/ansible/local.yaml
