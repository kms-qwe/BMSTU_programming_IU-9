#!/usr/bin/python3
from functools import lru_cache
def memoize(func):

    cache = {}    
    def memoized_func(*args):
  
        if args not in cache:
            cache[args] = func(*args)
        
        return cache[args]

    return memoized_func

def fib(n):
    if n < 0:
        return "Invalid argument"
    if n <= 2:
        return n 
    return fib(n - 1) + fib(n - 2)

fib = memoize(fib)
for i in range(10):
    print(fib(i))
