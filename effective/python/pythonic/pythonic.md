Effective Python Pythonic
---

```markdown
The Zen of Python, by Tim Peters

Beautiful is better than ugly.
Explicit is better than implicit.
Simple is better than complex.
Complex is better than complicated.
Flat is better than nested.
Sparse is better than dense.
Readability counts.
Special cases aren't special enough to break the rules.
Although practicality beats purity.
Errors should never pass silently.
Unless explicitly silenced.
In the face of ambiguity, refuse the temptation to guess.
There should be one-- and preferably only one --obvious way to do it.
Although that way may not be obvious at first unless you're Dutch.
Now is better than never.
Although never is often better than *right* now.
If the implementation is hard to explain, it's a bad idea.
If the implementation is easy to explain, it may be a good idea.
Namespaces are one honking great idea -- let's do more of those!
```

> The Zen of Python
> 每件事都应该有简单的做法，而且最好只有一种。
> 用直观、简洁而且容易看懂的方式来编写代码。

> [Python 官网][0]
> [Python 风格指南][1]
> [Pylint 自动检查受测试代码是否符合 PEP8 风格指南][2]
> [Slatkin 个人网站][3]
> [Effective Python Book][4]

### 里程碑

1. Python 2 已经在 2020/1/1 退场，停止更新维护。
   深度依赖 Python 2 代码库的开发者，可以考虑使用 `2to3 (Python 预装的工具)` 与 [six](https://six.readthedocs.io/) (社区包) 这样的工具，过渡到 Python 3。

### PEP 8 风格指南
#### 建议

##### 空白

- 用空格 space 缩进，而不要用制表符 tab
- 每行不超过 79 个字符
- 同一份文件中，函数与类之间用两个空行隔开
- 同一个类中，方法与方法之间用一个空行隔开
- 同一行的冒号和值间，应该加一个空格

##### 命名

- 函数、变量及属性，用小写字母拼写，之间用 `_` 相连
- 受保护的实例属性，用 `_` 开头
- 私有的实例属性，用 `__` 开头
- 类(包括异常)命名时，每个单词的首字母均大写，如 `CapWord`
- 模块级别的常量，所有字母都大写，之间用 `_` 相连
- 类中的实例方法，应该把第一个参数命名为 `self`，用来表示该对象本身
- 类方法的第一个参数，应该命名为 `cls`，用来表示类本本身

##### 表达式和语句

- (行内否定)否定词直接写在要否定的内容前：`if a is not b`
- 不要通过长度，判断序列是不是空，因为会把空值自动评估为 False，应该采用 `if not somelist`
- 不要通过长度，判断序列里有没有内容，因为会把非空值自动判定为 True，应该采用 `if somelist`
- 多行表达式，应该用括号括起来，而不要用 `\` 符号续行

#### 引入

- 引入模块时，总是应该使用绝对名称，如 `from bar import foo`
- 如果一定要用相对名称来编写 `import` 语句，应该明确地写成 `from . import foo`
- `import` 语句，应该划分成三部分：标准库中的模块、第三方模块、自己的模块

### 技巧

1. 通过二分法在有序的列表中搜索，让程序跑得更快
2. 只能通过关键字形式来指定的参数，把代码写得更加清晰易读
3. 使用星号表达式，拆分序列，减少出错率
4. 通过 `zip` 并行迭代多个列表，让代码更具 Python 风格

#### 简易操作

```python
$ python3 --version
Python 3.7.3
```

```python
$ python3
Python 3.7.3 (default, Mar 27 2019, 09:23:15)
[Clang 10.0.1 (clang-1001.0.46.3)] on darwin
Type "help", "copyright", "credits" or "license" for more information.
>>> import this
The Zen of Python, by Tim Peters

Beautiful is better than ugly.
Explicit is better than implicit.
Simple is better than complex.
Complex is better than complicated.
Flat is better than nested.
Sparse is better than dense.
Readability counts.
Special cases aren't special enough to break the rules.
Although practicality beats purity.
Errors should never pass silently.
Unless explicitly silenced.
In the face of ambiguity, refuse the temptation to guess.
There should be one-- and preferably only one --obvious way to do it.
Although that way may not be obvious at first unless you're Dutch.
Now is better than never.
Although never is often better than *right* now.
If the implementation is hard to explain, it's a bad idea.
If the implementation is easy to explain, it may be a good idea.
Namespaces are one honking great idea -- let's do more of those!

>>> import sys
>>> print(sys.version_info)
sys.version_info(major=3, minor=7, micro=3, releaselevel='final', serial=0)
>>> print(sys.version)
3.7.3 (default, Mar 27 2019, 09:23:15)
[Clang 10.0.1 (clang-1001.0.46.3)]
>>>
```

---
[0]: https://www.python.org/ "Python 官网"
[1]: https://www.python.org/dev/peps/pep-0008/ "PEP 8 - Style Guide for Python Code"
[2]: https://wwww.pylint.org/ "一款流行的 Python 源码静态分析工具"
[3]: https://onebigfluke.com "Slatkin 个人网站"
[4]: https://github.com/bslatkin/effectivepython
[5]: https://github.com/bslatkin/effectivepython/tree/master/example_code "Effective Python: Second Edition ExampleCode"
[6]: https://effectivepython.com/ "Effective Python: Second Edition"

> 2021/03/04 Katio