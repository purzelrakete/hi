"""
Test nn
"""

import numpy as np
import nn

np.random.seed(1)


def test_linear():
    l = nn.Linear(1, 2)
    got = l.forward([1, 2])
    assert got.shape == (1,)


def test_relu():
    l = nn.Relu()
    got = l.forward([-1, 1])
    assert list(got) == [0, 1]


def test_sigmoid():
    l = nn.Sigmoid()
    got = l.forward([0])
    assert list(got) == [0.5]


def test_sequential():
    X = np.random.randn(2, 5)
    y = np.random.randn(1, 5)

    predictor = (nn.Sequential()
                   .add(nn.Linear(5, 2))
                   .add(nn.Relu())
                   .add(nn.Linear(5, 1))
                   .add(nn.Sigmoid()))

    y_pred = predictor.forward(X)

    print(y_pred)
    print(y_pred.shape)
