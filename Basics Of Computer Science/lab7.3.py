#!/usr/bin/python3
from lab7mod import generate_random_string
import sys

if int(sys.argv[1]) <= 0 or int(sys.argv[2]) <= 0:
    print("Invalid arguments")
else:
    for i in range(int(sys.argv[2])):
        print(generate_random_string(int(sys.argv[1])))