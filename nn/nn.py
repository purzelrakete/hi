#!/usr/bin/env python3.6

"""
Feed Forward Neural Network in Numpy
"""

import numpy as np


class Linear:
    """A Linear Layer
    """

    def __init__(self, d1: int, d2: int):
        self.W = np.random.randn(d1, d2) / np.sqrt(2 / d2)
        self.b = 0

    def forward(self, X):
        self.dW = 1
        self.dX = 1
        self.db = 1
        return np.dot(self.W, X) + self.b

    def backward(self, dout):
        return dout * self.dX


class Relu:
    """A Relu Activation Function
    """

    def forward(self, X):
        return np.maximum(X, 0)

    def backward(self, dout):
        return dout * (self.input > 0)


class Sigmoid:
    """A Sigmoid Activation Function
    """

    def forward(self, Z):
        return self.sig(Z)

    def backward(self, dout):
        return dout * self.sig(self.input) * (1 - self.sig(self.input))

    def sig(self, z):
        return 1 / (1 + np.exp(np.multiply(-1, z)))


class CrossEntropy:
    """CrossEntropy Loss
    """

    def __init__(self, X, y):
        self.X = X
        self.y = y

    def forward(self):
        nll = -np.log(self.y * self.input + (1 - self.y) * (1 - self.input))
        return np.sum(nll, axis = 2)

    def backward(self, dout):
        return dout * 1


class Sequential:
    """A sequential Model, Keras style
    """

    def __init__(self):
        self.nodes = []

    def forward(self, X):
        return self.process(self.nodes, X)

    def backward(self):
        return self.process(reverse(self.nodes), 1)

    def process(self, nodes, inp):
        if len(nodes) == 0:
            return inp
        else:
            head, *tail = nodes
            return self.process(tail, head.forward(inp))

    def add(self, node):
        self.nodes.append(node)
        return self
