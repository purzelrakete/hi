"""
Data helpers.

"""

import numpy as np


def gaussian_clusters(nclasses, m, n, bounds = 10):
    "returns nclasses different clusters with labels X, Y"

    Xs = [gaussian_cluster(m, n, bounds) for cls in range(nclasses)]
    X = np.vstack(tuple(Xs))
    Ys = [np.repeat(cls, m) for cls in range(nclasses)]
    Y = np.hstack(tuple(Ys))

    return X.T, Y


def gaussian_cluster(m, n, bounds = 10):
    "produce a gaussian cluster with mean within +- bounds"

    return np.random.multivariate_normal(
        [np.random.uniform(-bounds, bounds) for _ in range(n)],
        cov(n, bounds),
        m)


def cov(dims, bounds = 10):
    """
    Covariance between two random variables is:

        cov(X, Y) = E[(X - E[X] * (Y - E[Y])]

    The covariance Matrix entry (i, j) is cov(Ri, Rj) given a set of m
    random variables. So in 2 dimensions, we have a 2x2 matrix, the diagonals
    describe the marginal variance, the off diagonals the joint variance
    between Ri, Rj. Covariance Matrices are symmetric, eg X' = X.

    The covariance matrix has to be x'Ax > 0, eg positive semidefinite.
    """

    # keep going until we find a matrix with all positive eigenvalues.
    # symmetric matrices with this property are positive semidefinite.
    ret = np.ones((dims, dims)) * -1
    while np.any(np.linalg.eig(ret)[0] < 0):
        ret = np.zeros((dims, dims))

        # variance along dimension i
        for i in range(dims):
            ret[i, i] = np.random.uniform(-bounds, bounds)


        # covariance of variables i, j
        for i in range(dims):
            for j in range(i + 1, dims):
                val = np.random.uniform(-bounds, bounds)
                ret[i, j] = val
                ret[j, i] = val

    return ret
