"""
Logistic regression in numpy.

"""

import numpy as np


def train(X, Y):
    "train and return a model"

    n, m = X.shape

    def optimize(iterations = 1000, learning_rate = 0.001):
        W, b = np.zeros((n, 1)), 0

        for i in xrange(iterations):
            dW = 0
            db = 0

            W = W - learning_rate * dW
            b = b - learning_rate * db

        return Model(W, b)

    model = optimize()
    loss = model.nll(X, Y)

    return model


class Model(object):
    "Logistic Regression model"

    def __init__(self, W, b):
        self.W = W
        self.b = b

    def predict(self, X, threshold = 0.5):
        A = sigmoid(np.dot(self.W.T, X) + self.b)
        classes = [1 if p > threshold else 0 for p in A]

        return classes, A

    def boundary(self, x):
        "Decision boundary in the R2 case."

        assert self.W.shape[0] == 2
        return self.b + self.W[0] * x / self.W[1]

    def nll(self, X, Y):
        "Negative Log Likelihood of given X, Y under this model"

        n, m = X.shape
        A = sigmoid(np.dot(self.W.T, X) + self.b)
        logprobs = np.log(A**Y) * (1-A)**np.log(1-Y)

        return -sum(logprobs) / m


#
# Maths
#

def sigmoid(x):
    return 1 / (1 + np.exp(-x))
