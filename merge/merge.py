#!/usr/bin/env python
# -*- coding: utf-8 -*-
#
# Let's sort random stuff!

import heapq


def it(ll):
    """A a sorting generator for the given list of sorted lists.

    Args:
        - a list of lists. lists may contain any sortable data type. The
          inner lists should already be sorted.

    Returns:
        - a generator which will yield elements in sorted order, lowest to
          highest.
    """

    # remove empty lists first
    ll = [l for l in ll if l]
    if not ll:
        return

    # initial state
    k = len(ll)
    pointers = k * [0]
    lens = [len(l) for l in ll]
    heap = list(zip([l[0] for l in ll], range(k)))
    heapq.heapify(heap)

    for _ in range(sum(lens)):
        lowest_val, i = heapq.heappop(heap)
        if pointers[i] < lens[i] - 1:
            pointers[i] += 1
            heapq.heappush(heap, (ll[i][pointers[i]], i))

        yield lowest_val


def it2(ll):
    """A a iterator for the given list of sorted lists.

    Args:
        - a list of lists. lists may contain any sortable data type. The
          inner lists should already be sorted.

    Returns:
        - an iterator which will yield elements in sorted order, lowest to
          highest.
    """

    return heapq.merge(*ll)
