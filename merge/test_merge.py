#!/usr/bin/env python
# -*- coding: utf-8 -*-
#
# Test merge

import itertools

from hypothesis import given
from hypothesis.strategies import lists, integers

import merge


@given(lists(lists(integers())))
def test_it(ll):
    ll = [sorted(l) for l in ll]
    assert list(merge.it(ll)) == sorted(list(itertools.chain.from_iterable(ll)))


@given(lists(lists(integers())))
def test_it2(ll):
    ll = [sorted(l) for l in ll]
    assert list(merge.it2(ll)) == sorted(list(itertools.chain.from_iterable(ll)))
