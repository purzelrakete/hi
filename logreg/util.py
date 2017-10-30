# -*- coding: utf-8 -*-

"""
Maths

"""

import numpy as np


def sigmoid(x):
    return 1 / (1 + np.exp(-x))
