IPython Interactive Computing
---

IPython provides a rich architecture for interactive computing with:

- a powerful interactive shell
- a kernel for jupyter
- support for interactive data visualization and use of GUI toolkits
- flexible, embeddable interpreters to load into your own projects
- easy to use, high performance tools for parallel computing

### Skips

- Usage

```shell
pip3 install ipython
```

```ipython
ipython

Python 3.7.3 (default, Mar 27 2019, 09:23:15)
Type 'copyright', 'credits' or 'license' for more information
IPython 7.32.0 -- An enhanced Interactive Python. Type '?' for help.

In [1]: ?

# print var detail info
# ?? 还可以查看函数或模块对象的源码
In [2]: object?
In [3]: list?
In [4]: list??

# 历史命令：hist, history
In [5]: hist
In [6]: history

# tab 自动补全
# ！shell_command 执行 shell 命令
In [7]: ! ping www.baidu.com

# 魔法命令：line magics - %page 表示魔法只在本行有效，cell magics %%page 表示魔法在整个单元有效
# %run 运行脚本：%run 路径+文件名
# %timeit 测量单行代码的运行时间: %timeit [i for i in range(1000)]
# %%timeit 测量整个单元代码运行时间: 
%%timeit
obj = []
for i range(100):
    obj.append(i)

%env 显示环境变量
%pwd 显示当前工作目录
%ls path 显示特定目录下的内容
%cd
_ 打印前一个输出结果，它是一个变量，实时更新的
__ 获取倒数第二个输出结果
___ 获取倒数第三个输出结果
在语句后面加上 ;，不显示输出结果
%pdef 打印类、函数的构造信息
%pdoc 打印对象文档字符串
import numpy as np
%pdoc np

%conda 安装python第三方库
# 在 notebook 中绘制图像时，将图表直接嵌入到 notebook 中
%matplotlib inline
import matplotlib.pyplot as plt
%matplotlib inline

obj = range(5)
plt.plot(obj)

%pylab 使numpy和matplotlib中的科学计算功能生效
%pylab
%matplotlib inline
x=pylab.linespace(-10.,10.,1000)
plot(x,sin(x))

%quickref 查看 IPython 特定语法和魔法命令参考

*? 模糊查询方法名及属性 -> 
import pandas as pd
pd.*Da*?

%debug 从最新的异常跟踪的底部进入交互式调试器
在ipdb调试模式下能访问所有的本地变量和整个栈回溯。使用u和d向上和向下访问栈，使用q退出调试器。在调试器中输入?可以查看所有的可用命令列表。
%pdb 支持对所有的异常进行调试；需要事先启动%pdb命令，之后对每一个异常都会进行调试。
%run -d test.py 用于对脚本进行调试
%pycat filename 语法高亮显示一个 python 文件
%load test.py 将脚本代码加载到当前单元
%notebook path 导出当前 notebook 内容到指定 ipynb 文件中

%precision 设置浮点数精度，可添加具体参数，无参数则默认精度
from math import pi
pi
%precision 3
pi

%xdel 用于删除变量，并尝试清楚其在IPython中的对象上的一切引用
%who  用于显示当前所有变量，也可以指定显示变量的类型
%who
%who int
%whos  用于显示当前所有变量，也可以指定显示变量的类型，但提供的信息更加丰富
%save path n1 n2.. 用于将指定cell代码保存到指定的py文件中
%reset -f 用于删除定义的所有变量，如果不指定参数-f，则需要确认后再重置
%rerun 用于执行之前的代码，可以指定历史代码行，默认最后一行
%%latex 用于将LaTeX语句渲染为公式，LaTeX是一种基于ΤΕΧ的排版系统
%%html 用于将单元格渲染为HTML输出
%%js 用于运行含有JavaScript代码的cell
%%markdown 用于将markdown文本渲染为可视化输出
%%writefile path 用于将单元格内容写入到指定文件中，文件格式可为txt、py等
%bookmark 能够保存常用目录的别名，以便实现快速跳转，书签能够持久化保存
%paste 能够直接执行剪切板中的python代码块
%magic 获取所有魔法命令及其用法

In 对象是一个列表，按照顺序记录所有的命令。
Out 对象不是一个列表，而是一个字典，它将输入数字映射到相应的输出（如果有的话）

%xmode 用于控制异常输出的模式
```

[0]: https://ipython.org/index.html "IPython 官网"
[1]: https://pypi.org/project/ipython/#history "the release history for ipython"