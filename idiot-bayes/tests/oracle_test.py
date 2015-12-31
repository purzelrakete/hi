# -*- coding: utf-8 -*-

"""
    oracle tests
    ~~~~~~~~~~~~

    :copyright: (c) 2015 Rany Keddo
    :license: BSD, see LICENSE for more details.
"""

from polyglot import oracle


def test_oracle():
    file_name = 'dataset/go/sieve.go'
    p_go = oracle.predict(file_name)

    assert p_go > 0.0
    assert p_go < 1.0
