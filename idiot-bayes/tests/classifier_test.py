# -*- coding: utf-8 -*-

"""
    classifier tests
    ~~~~~~~~~~~~~~~~

    :copyright: (c) 2015 Rany Keddo
    :license: BSD, see LICENSE for more details.
"""

from polyglot import classifier

def test_notgo_dataset_loads():
    ds = classifier.dataset('dataset/notgo')
    assert len(ds) > 100
    assert len(ds[0]) > 0

def test_go_dataset_loads():
    ds = classifier.dataset('dataset/go')
    assert len(ds) > 100
    assert len(ds[0]) > 0

def test_train():
    model = classifier.train()
    p_go_given_token = model.p_go_given_token('range')

    # counts look ok
    assert model.nall_docs > 400
    assert model.nnotgo_docs > 150
    assert model.ngo_docs > 150

    # probabilities within reasonable bounds
    assert p_go_given_token > 0.5
    assert p_go_given_token < 1.0
