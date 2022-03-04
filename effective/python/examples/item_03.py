#!/usr/bin/env PYTHONHASHSEED=1234 python3

import random

random.seed(1234)

import logging
from pprint import pprint
from sys import stdout as STDOUT

# Write all output to a temporary directory
import atexit
import gc
import io
import os
import tempfile

TEST_DIR = tempfile.TemporaryDirectory()
atexit.register(TEST_DIR.cleanup)

# Make sure Windows processes exit cleanly
OLD_CWD = os.getcwd()
atexit.register(lambda: os.chdir(OLD_CWD))
os.chdir(TEST_DIR.name)


def close_open_files():
    everything = gc.get_objects()
    for obj in everything:
        if isinstance(obj, io.IOBase):
            obj.close()


atexit.register(close_open_files)


# 字符序列实例，不一定非要用某一种固定的方案编码成二进制数据，要把二进制数据转换成 Unicode 数据，必须调用 bytes decode 方法，
# 而要把 Unicode 数据转换成二进制数据，必须调用 str encode 方法
# Example 字符序列：bytes
# bytes 实例包含的是原始数据，即 8 位的无符号值(ASCII 编码标准)
a = b'\x65llo'
print(list(a))
print(a)


# Example 字符序列：str
# str 实例包含的是 Unicode 码点(code point 代码点)，这些码点与文本字符相对应
a = 'a\u0300 propos'
print(list(a))
print(a)


# Example 3