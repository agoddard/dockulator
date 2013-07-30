#!/usr/bin/env python
import sys
import operator

operator_dict = {
    "+": operator.add,
    "-": operator.sub,
    "*": operator.mul,
    "/": operator.div,
}

def main():
    left, op, right = sys.argv[1].split()
    answer = operator_dict[op](int(left), int(right))

    print answer
    

if __name__ == "__main__":
    main()

