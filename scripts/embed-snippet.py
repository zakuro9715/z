#!/usr/bin/env python3
import sys
from glob import glob

usage = '''
md-embed-snippet template placeholder snippets

Args:
    template    - template name. e.g. README.template.md
    placeholder - text to replace. e.g. __examples__
    snippets    - files to embed e.g. examples/*
'''


def expand_examples(arg):
    return sum([glob(v) for v in arg], [])


def snippet(filepath):
    return "```{}\n{}```".format(
        filepath,
        open(filepath).read(),
    )


def main():
    args = sys.argv[1:]
    if len(args) < 3:
        print(usage)
        return
    template = open(args[0]).read()
    placeholder = args[1]
    examples = "\n\n".join([snippet(f) for f in expand_examples(args[1:])])
    text = template.replace(placeholder, examples)
    print(text, end="")


if __name__ == '__main__':
    main()
