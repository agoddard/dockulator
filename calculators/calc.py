import sys
import os
import operator

operator_dict = {
    "+": operator.add,
    "-": operator.sub,
    "*": operator.mul,
    "/": operator.div,
}

def main():
    left, op, right = sys.argv[1].split()
    print(operator_dict[op](int(left), int(right)))

    print(sys.argv[2])
    print(os.popen('hostname').read()),

if __name__ == "__main__":
    main()

