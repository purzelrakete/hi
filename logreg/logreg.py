# -*- coding: utf-8 -*-

"""
Logistic regression in numpy.

"""

import numpy as np

import util


def train(X, Y, iterations = 1000, learning_rate = 0.001, momentum = 0.9):
    "train and return a model"

    n, m = X.shape

    def optimize():

        # W ∈ Rn, b ∈ R1
        W, b = np.zeros((n, 1)), 0
        losses = []

        for i in xrange(iterations):

            # activations
            A = util.sigmoid(np.dot(W.T, X) + b)

            # dJ/dW ∈ Rn
            dW = np.dot(X, (A - Y).T) / m

            # dJ/db ∈ R1
            db = np.sum(A - Y) / m

            W = W - learning_rate * dW
            b = b - learning_rate * db

            losses.append(Model(W, b).nll(X, Y))

        return Model(W, b), losses

    model, losses = optimize()

    return model, losses


class Model(object):
    "Logistic Regression model"

    def __init__(self, W, b):
        self.W = W
        self.b = b

    def predict(self, X, threshold = 0.5):
        n, m = X.shape
        A = util.sigmoid(np.dot(self.W.T, X) + self.b)

        assert A.shape == (1, m)
        classes = np.where(A > 0.5, 1, 0)

        return classes, A

    def boundary(self, x):
        "Decision boundary in the R2 case. This is orthogonal to W."

        assert self.W.shape == (2, 1)
        return (-1.0 * self.W[0] * x - self.b) / self.W[1]

    def nll(self, X, Y):
        "Negative Log Likelihood of given X, Y under this model"

        n, m = X.shape
        A = util.sigmoid(np.dot(self.W.T, X) + self.b)

        # A invariants
        assert A.shape == (1, m)
        assert np.sum(A >= 0) == m
        assert np.sum(A <= 1) == m

        # Y invariants
        assert Y.shape == (1, m)
        assert np.sum(Y == 1) + np.sum(Y == 0) == m

        logprobs = Y * np.log(A) + (1 - Y) * np.log(1 - A)

        # logprob invariants
        assert logprobs.shape == (1, m)
        assert np.sum(logprobs <= 0) == m

        return np.squeeze(-np.sum(logprobs, axis = 1, keepdims = True) / m)
