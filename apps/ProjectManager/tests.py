# -*- coding:utf-8 -*-

name = '付宗飞'

print(name[::-1])


import sys
import io
sys.stdout = io.TextIOWrapper(sys.stdout.buffer, encoding='utf-8')