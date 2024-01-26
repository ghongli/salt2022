# MacOS

## Tools

### [MacPorts][1]

[查找安装包][2]

```
sudo port install git
sudo port install wget

# net utils: telnet
#sudo port install p5.30-net-telnet
sudo port install inetutils
gtelnet --help
```

[1]: https://www.macports.org/ "MacPorts Home"
[2]: https://ports.macports.org/ "包管理"

### [Brew][11]

[更多文档][12]
[brew Shell Completion](https://docs.brew.sh/Shell-Completion)

```
# install
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

sudo brew install git
```

[11]: https://brew.sh/zh-cn/
[12]: https://docs.brew.sh/ "更多文档"
