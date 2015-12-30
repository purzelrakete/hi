# s -*- coding: utf-8 -*-

"""
    classifier
    ~~~~~~~~~~

    Naive Bayes classifer written during a train trip to Berlin.

    Classify programm code into { go, notgo }.

    :copyright: (c) 2015 Rany Keddo
    :license: BSD, see LICENSE for more details.
"""

import math
import operator
import os

def classify(code, model):
    """
    returns the probabilty that the given code is go.
    """

    p_token = map(model.p_go_given_token, tokenize(code))
    return math.pow(10, sum(map(math.log10, p_token)))


def train():
    """
    trains the classifier on a fixed dataset.
    """

    return Model(
        go_docs = dataset('dataset/go'),
        notgo_docs = dataset('dataset/notgo')
    )


def dataset(dirname):
    """
    reads a set of files from a given directory, converting each to
    a document.
    """

    docs = [document(dirname + '/' + filename) for filename in os.listdir(dirname)]
    return [doc for doc in docs if len(doc) > 0]


def document(filename):
    """
    convert an on disk program into a set of tokens.
    """

    with open(filename, 'r') as f:
        return tokenize(f.read())


def tokenize(code):
    """
    tokenize program code. most stupid possible tokenizer splits on ws.

    returns a set of uniqe tokens
    """

    return [token for token in set(code.split(' ')) if token != '']


def with_token(token, documents):
    """returns the number of documents containing the token. """

    return len([d for d in documents if token in d])


# TODO(rk): review this properly, eg what MLE?
def smoothed_mle(a, b):
    """
    laplacian smoothed maximum likelihood
    """

    # the number of classes in our classifier
    k = 2

    return float(a + 1) / float(b + k)


class Model:
    """trained model. """

    def __init__(self, go_docs, notgo_docs):
        self.go_docs = go_docs
        self.notgo_docs = notgo_docs

        # counts
        self.ngo_docs = len(self.go_docs)
        self.nnotgo_docs = len(self.notgo_docs)
        self.nall_docs = self.ngo_docs + self.nnotgo_docs

        # priors
        self.p_go = smoothed_mle(self.ngo_docs, self.nall_docs)


    def p_go_given_token(self, token):
        """
        probability this is go given the token. applies bayes theorem to
        do inference.
        """

        # counts
        ngo_docs_with_token = with_token(token, self.go_docs)
        nnotgo_docs_with_token = with_token(token, self.notgo_docs)
        nall_docs_with_token = ngo_docs_with_token + nnotgo_docs_with_token

        # probabilities
        p_token_given_go = smoothed_mle(ngo_docs_with_token, self.ngo_docs)
        p_token = smoothed_mle(nall_docs_with_token, self.nall_docs)

        return p_token_given_go * self.p_go / p_token
