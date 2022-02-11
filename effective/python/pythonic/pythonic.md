Effective Python Pythonic
---

> The Zen of Python
> 每件事都应该有简单的做法，而且最好只有一种。
> 用直观、简洁而且容易看懂的方式来编写代码。

> [Python 官网][0]
> [Python 风格指南][1]
> [Pylint 自动检查受测试代码是否符合 PEP8 风格指南][2]
> [Slatkin 个人网站][3]
> [Effective Python Book][4]

### 里程碑
1. Python 2 已经在 2020/1/1 退场。

### PEP 8 风格指南
#### 建议



##### 表达式和语句

- (行内否定)否定词直接写在要否定的内容前：`if a is not b`
- 不要通过长度，判断序列是不是空，因为会把空值自动评估为 False，应该采用 `if not somelist`
- 不要通过长度，判断序列里有没有内容，因为会把非空值自动判定为 True，应该采用 `if somelist`
- 多行表达式，应该用括号括起来，而不要用 `\` 符号续行

### 技巧

1. 通过二分法在有序的列表中搜索，让程序跑得更快
2. 只能通过关键字形式来指定的参数，把代码写得更加清晰易读
3. 使用星号表达式，拆分序列，减少出错率
4. 通过 `zip` 并行迭代多个列表，让代码更具 Python 风格

---
[0]: https://www.python.org/ "Python 官网"
[1]: https://www.python.org/dev/peps/pep-0008/ "PEP 8 - Style Guide for Python Code"
[2]: https://wwww.pylint.org/ "一款流行的 Python 源码静态分析工具"
[3]: https://onebigfluke.com "Slatkin 个人网站"
[4]: https://github.com/bslatkin/effectivepython

> 2021/03/04 Katio