SHELL:=/usr/bin/env bash
.DEFAULT_GOAL:=all

MM_SHELL ?= bash -c
MAKEFLAGS += --no-print-directory

# ansible-playbook params:
# -v, --verbose
#	 verbose mode (-vvv for more, -vvvv to enable connection debugging)
# -D, --diff
#	 when changing (small) files and templates, show the differences in those files
# -e EXTRA_VARS, --extra-vars EXTRA_VARS
#	 set additional variables as key=value or YAML/JSON, if filename prepend with @
# -i INVENTORY, --inventory INVENTORY
#    specify inventory host path or comma separated host list
# -l 限制需要部署的节点范围，如 `frontend`, `apiserver`, `backend`, `license` 等
# -t 根据相应的值，仅仅运行选中的任务，类 k8s selected 策略，如 `install`, `force`,`start`, `stop`, `restart`, `license up` 等
# -k, --ask-pass
#	  ask for connection password
# -K, --ask-become-pass
#     ask for privilege escalation password

.PHONY: syntax-check # 检查 playbook 脚本的语法问题
syntax-check:
	@ansible-playbook -D -v --syntax-check site.yml

.PHONY: debug # 执行 playbook 脚本，但会输出更详细的信息
debug: syntax-check
	@ansible-playbook -D -vvvv -i hosts site.yml

.PHONY: deploy # 执行 playbook 脚本
deploy: syntax-check
	@ansible-playbook -D -v -i hosts site.yml

START_TASK ?= ''

.PHONY: startTask # 在某个任务开始执行 playbook 脚本.
startTask: syntax-check
	@ansible-playbook -D -v site.yml --start-at='${START_TASK}'

.PHONY: step # one-step-at-a-time: confirm each task before running.
step: syntax-check
	@ansible-playbook -D -v site.yml --step

SETUP_FILTER ?= '*cpu*'

.PHONY: setupFilter # Gathers facts about localhost by some filter, e.g. '*mb*', '*cpu*'
setupFilter:
	@ansible localhost -m setup -a 'filter=${SETUP_FILTER}'

.PHONY: help # print this help message.
 help:
	@grep -E '^\.PHONY: [a-zA-Z_-]+ .*?# .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = "(: |#)"}; {printf "%-30s %s\n", $$2, $$3}'