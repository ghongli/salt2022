Ansible
---

> Alan Kay 格言 “简单的事情应该保持简单，复杂的事情应该做到可能”。

Ansible 引入了一种特定领域的语言(DSL domain-specific language)描述服务器的状态。Ansible 的设计初衷就是在若干服务器上从零开始执行所有必需的配置与操作。
Ansible 使用极简的模型来实现对各种操作按照所需顺序执行的控制。

对于 Ansible 来说，脚本被称作 playbook。playbook 描述了需要配置的主机(远程服务器)和主机上需要按顺序运行的任务列表，依次进行任务列表中的每一个任务。
需要格外注意的是：
- 对于每一个任务，Ansible 都是在所有主机之间并行执行的。
- 在开始下一个任务之前，Ansible 会等待所有主机都完成上一个任务。
- Ansible 会按照指定的顺序来运行任务。

Ansible 不需要使用者去学习一套新的用于屏蔽不同操作系统差异的抽象。这使得 Ansible 需要学习的东西更少，在开始使用它之前，几乎不需要其他必备的知识。
如果真想有这层抽象，可以编写自己的 Ansible playbook 时，实现针对不同操作系统的远程服务器，运行不同的操作。
`但实际中，尽量避免这么做，更专注于编写用于某个特定的操作系统(e.g Ubuntu)的 playbook 脚本。`

Ansible 使用 SSH 的多路复用技术(multiplexing)来优化性能。

Ansible 使用 YAML 文件格式和 Jinja2 模板语言，所以对于使用 Ansible，需要学习一些 YAML、Jinja2 的知识，但它们很容易上手。

Ansible 只能管理那些它明确了解的服务器，即在 inventory 文件中指定的服务器信息。通过使用 `ansible.cfg` 文件来避免冗长的 inventory 文件。
`ansible.cfg` 文件可以设定一些默认值，简化配置，不需要将同样的内容键入很多遍。

Ansible 按照如下的位置和顺序查找 `ansible.cfg` 文件：
1. `ANSIBLE_CONFIG` 环境变量所指定的文件
2. `./ansible.cfg` (当前目录下的 ansible.cfg)
3. `~/.ansible.cfg` (主目录下的 .ansible.cfg)
4. `/etc/ansible/ansible.cfg`

Ansible 使用 `/etc/ansible/host` 为 inventory 文件的默认位置，实际上 inventory 文件和 playbook 一起进行版本控制。

### 特色

#### 易读的语法

脚本 playbook 的语法是基于 YAML 构建的，记录了部署所须使用的命令，只不过这些指令永远不会过期，因为它们同时也是直接执行的代码。

#### 远程主机无需安装依赖

需要被 Ansible 管理的服务器，需要安装 OenSSH 和 Python，不再需要预装任何代理程序或其他软件了；控制服务器，需要安装 Python。

#### 基于推送模式，是一种无 agent 模型

Ansible 默认采用的是基于推送的模式，配置变更步骤如下：
- 在 playbook 中进行变更
- 运行新的 playbook
- 连接到服务器，并执行那些改变服务器状态的模块

优点：直接由脚本执行者，控制变更在服务器上发生的时间。

> 如果更喜欢拉取模式，Ansible 可以使用 `ansible-pull` 工具实现。

#### 声明式的模块

模块声明它可以用于描述希望服务器所处于的状态。如使用 `user 用户模块`，保证拥有一个名有 deploy 且所属组为 web 的账号：`user: name=deploy group=web`。

模块是幂等的。如果用户 deploy 不存在，Ansible 就创建它，否则不会做任何事。`幂等性，意味着向同一台服务器多次执行同一个 Ansible playbook 是安全的。`

Ansible 并没有需要多次运行来配置服务器的设计。相对的，Ansible 的模块实现的行为是：只需要运行 playbook 一次，即可将每台服务器都配置为期望的状态。

Ansible 社区，复用的基本单元是模块。由于模块的功能范围非常小，并且可以只针对特定的操作系统，所以易于实现定义明确且分享的模块。

### 在线资源

> [Ansible][0], [Ansible Docs][10], [Ansible Examples][4], [Ansible Repo][1], [Ansible Releases][5] | [Tarballs of releases][16], [Ansible 角色仓库][7]
> [Ansible 开发组的 Google Group][3], [Ansible 5 模块索引][13], [Ansible 2.9 模块索引][14]
> [Ansible Book Repo][20], [简明参考手册][21]

[0]: https://www.ansible.com/
[1]: https://github.com/ansible
[2]: https://github.com/ansible/ansible
[3]: https://groups.google.com/g/ansible-project "Ansible 开发组的 Google Group"
[4]: https://github.com/ansible/ansible-examples
[5]: https://github.com/ansible/ansible/releases
[6]: https://pypi.org/project/ansible/ "PyPI's ansible package page"
[7]: https://galaxy.ansible.com/ "Ansible 角色仓库(Ansible Galaxy)"

[10]: https://docs.ansible.com/ "官方文档"
[11]: https://docs.ansible.com/ansible/latest/index.html "ansible latest index"
[12]: https://docs.ansible.com/ansible/latest/user_guide/basic_concepts.html#basic-concepts "ansible concepts"
[13]: https://docs.ansible.com/ansible/latest/collections/index_module.html "模块索引"
[14]: https://docs.ansible.com/ansible/2.9/modules/modules_by_category.html "2.9 模块索引"
[15]: https://docs.ansible.com/ansible/latest/installation_guide/intro_installation.html "Installing Ansible"
[16]: https://releases.ansible.com/ansible/ "Tarballs of releases"

[20]: https://github.com/ansiblebook/ansiblebook "Ansible: Up and Running, Second Edition"
[21]: https://github.com/lorin/ansible-quickref "简明参考手册"
[22]: http://www.slideshare.net/JesseKeating/ansiblefest-rax "使用 Ansible 管理可伸缩的公有云"

[30]: https://groups.google.com/g/ansible-project/c/WpRblldA2PQ/m/lYDpFjBXDlsJ "幂等性、收敛性以及那些我们说烂了的 “高达上” 术语"

[40]: http://www.vagrantup.com "Vagrant 是一个优秀的开源虚拟机管理工具，内置支持使用 Ansible 部署虚机"
[41]: http://www.virtualbox.org "Virtualbox 是 Vagrant 依赖项"
