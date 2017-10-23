"""
Test logistic regression

"""

import numpy as np

import logreg


def model(n):
    w, b = np.random.rand(n, 1), np.random.randn()
    return logreg.Model(w, b)


def data(m, n):
    return np.random.rand(n, m)


def test_prediction():
    m, n = 1, 2
    for _ in xrange(1000):
        Y, p = model(n).predict(data(m, n))
        assert len(p) == 1
        assert len(Y) == 1
        assert p[0] > 0
        assert p[0] < 1.


def test_training():
    m, n = 1, 2
    X = data(m, n)
    y = np.random.randint(0, 1)
    model = logreg.train(X, y)
    classes, _ = model.predict(X)

    assert True
