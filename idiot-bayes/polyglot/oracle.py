# -*- coding: utf-8 -*-

"""
    oracle
    ~~~~~~

    drives the classifier

    :copyright: (c) 2015 Rany Keddo
    :license: BSD, see LICENSE for more details.
"""

from polyglot import classifier
import argparse
import os
import sys


def main():
    """main method"""

    config = cfg()
    validate(config)
    p_go = predict(config['file_name'])
    print 'p_go is %f' % (p_go * 100)


def predict(file_name):
    """predict for the given file name."""

    model = classifier.train()

    with open(file_name, 'r') as f:
        return classifier.classify(f.read(), model)


def validate(config):
    """input validations"""

    if not os.path.exists(config["file_name"]):
        fail('cannot find ' + config["file_name"])


def cfg():
    """load the configuration"""

    parser = argparse.ArgumentParser(description='classify a file into notgo or go')
    parser.add_argument('--file-name', required=True, help='the file name')
    args = vars(parser.parse_args())

    return args


def fail(msg):
    """print an error message and exit"""

    sys.stderr.write(msg + '\n')
    sys.exit(1)
