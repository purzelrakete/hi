"""
Logistic regression in numpy.

"""

import numpy as np


def train(X, y):
    "train a model; return the trained thing"

    n, m = X.shape
    w, b = np.random.rand(n, 1), np.random.randn()

    return Model(w, b)


class Model(object):
    "Logistic Regression model"

    def __init__(self, w, b):
        self.w = w
        self.b = b

    def predict(self, X, threshold = 0.5):
        def sigmoid(x): return 1 / (1 + np.exp(-x))
        probabilities = sigmoid(np.dot(self.w.T, X) + self.b)
        classes = [1 if p > threshold else 0 for p in probabilities]

        return classes, probabilities

    def boundary(self, x):
        "Decision boundary in the R2 case. Will error if != R2"

        assert self.w.shape[0] == 2
        return self.b + self.w[0] * x / self.w[1]


#
# Random data
#

def gaussian_clusters(nclasses, m, n, bounds = 10):
    "returns nclasses different clusters with labels X, y"

    X = np.vstack(tuple([gaussian_cluster(m, n, bounds) for cls in range(nclasses)]))
    y = np.hstack(tuple([np.repeat(cls, m) for cls in range(nclasses)]))

    return X.T, y


def gaussian_cluster(m, n, bounds = 10):
    "produce a gaussian cluster with mean within +- bounds"

    mean = [np.random.uniform(-bounds, bounds) for _ in range(n)]
    cov = rnd_cov(n)

    return np.random.multivariate_normal(mean, cov, m)


def rnd_cov(n):
    "only works in R2"

    assert n == 2
    cov = np.eye(n)
    cov[0, 1] = np.random.uniform(-1, 1)
    cov[1, 0] = cov[0, 1]

    return cov
