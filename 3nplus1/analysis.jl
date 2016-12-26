using Gadfly
using DataFrames
using Distributions
using Base.Test

include("collatz.jl")
include("tree.jl")
include("viz.jl")

# color scheme
theme = Theme(
  background_color = colorant"black",
  default_point_size = 0.5mm)

# Terminology
# -----------
#
# - the reverse tree starts at 1. right branches are always n / 2 predecessors,
#   and left branches are always 3n + 1 predecessors. see the png.
# - a series is in the form p2^n, where p is odd.
# - a series has a single root, p, which is always an odd number.
# - a series can be found in the reverse tree by following down right from an
#   odd root.
# - a series may have left roots. these are the left branches of that series.
# - the left roots of a series are given as (roots, parent), where parent is
#   the parent series, identified by the root p.
#

# Long term behaviour
# -------------------------
#
# let's have a look at a single collatz sequence starting at 2016. these
# sequences are also known as hailstone sequences, because of the way hailstones
# are formed in the clouds: they keep drifting up and back down, until finally
# they accumulate enough mass and drop down to earth, ie 1. here are the ups and
# downs for 2016:
df = DataFrame(y = f(2016))
plot(df, y = :y, Geom.line, theme)

# once again  plotted with points, which reveals apparent structure in the
# peaks:
plot(df, y = :y, Geom.point, theme)

# now let's have a look at the sequence lengths, eg the number of steps it takes
# to get to one, given starting numbers up to 10k:
@time df = df_f_len(2:10_000)
plot(df, y = :y, Geom.point, theme)

# more apparent structure. let's take a look at the histogram of sequence
# lengths to get a feeling for the distribution:
@time df = df_f_len(2:10_000_000)
plot(df, x = :y, Geom.histogram, theme)

# let's have a look at the density:
plot(df, x = :y, Geom.density, theme)

# the histogram looks a little bit bimodal. let's have a closer look, maybe
# binning is hiding some information:
plot(df, x = :y, Geom.histogram(bincount = 1000), theme)

# looks like some structure. at any rate, it's a right skewed distribution.
# maybe it could be approaching a gamma distribution. let's try to fit it:
gamma = fit_mle(Gamma, df[:y])
gamma_pdf = pdf(gamma, 1:700)
plot(layer(df, x = :y, Geom.density, theme, order = 1),
     layer(df, y = gamma_pdf, Geom.line, Theme(default_color = colorant"red"), order = 2),
     theme)

# we could now look at the KL divergence between the data and the fit, and try
# to find the best distribution. I WILL NOT.

# ok, but surely sequence lengths will grow as we get extremely large numbers,
# since the mean distance to the root should increase. let's look at how the
# distribution changes as we get up  higher:
@time df = df_f_len_many([2:1_000, 2:10_000, 2:100_000, 2:1_000_000, 2:10_000_000])
plot(df, x = :y, color = :color, Geom.density, theme)

# Left Roots
# ----------
#
# left roots starting from p = 1 series, log scale
df = DataFrame(y = series_left_roots(1, 50))
plot(df, y = :y, Geom.line, Scale.y_log10, theme)

# these left roots are successively * 4 + 1. looking at this in Base
# 4, we can see it is the same as successively shifting left and adding
# 1, eg creating a string of 1's. multiplying these by 3 and adding 1
# yields 4^n. So the closed from for this is (4^n-1)/3. let's confirm:
@assert [log2(n * 3 + 1) / 2 for n in df[:y]] == [2.0:25.0;]

# what about all of those new p's? what sorts of left
# branching behaviour do these series have? eg p = 5, 21, 85, 341 etc. 5:
df = DataFrame(y = series_left_roots(5, 50))
plot(df, y = :y, Geom.line, Scale.y_log10, theme)

# this one also looks like * 4 + 1.  however this one starts at 3. so shifting
# 1 left and adding 1 will produce 3111etc. no we can't just multiply by 3
# to make this 333etc. instead we have to ((4^n - 1) / 3) + 2 * 4^n. let's
# subtract the two series to find the remaining quadratic term:
@assert (series_left_roots(1, 52) - series_left_roots(5, 50)) == [2 * 4^n for n in 0:24]

# interestingly, series p = 21 may not have left branches:
length(series_left_roots(21, 500))

# 85?
df = DataFrame(y = series_left_roots(85, 50))
plot(df, y = :y, Geom.line, Scale.y_log10, theme)

# again, * 4 + 1. 113 in base 4 is 1301:
#
# step 0 -> 1301    = ((4^0 - 1) / 3) + 113 * 4^0
# step 1 -> 13011   = ((4^1 - 1) / 3) + 113 * 4^1
# step 2 -> 130111  = ((4^2 - 1) / 3) + 113 * 4^2
# step 3 -> 1301111 = ((4^3 - 1) / 3) + 113 * 4^3 = 7253

# closed formula as above:
closed(y, n) = BigInt(((BigInt(4)^n - 1) / 3) + y * BigInt(4)^n)

# let's try all secondary series starting off the primary (1) series and see if
# the same closed form works for their tertiary series.
roots = [(x, series_left_roots(x, 50)) for x in series_left_roots(1, 50)]
filtered = filter(x -> !isempty(x[2]), roots)

# it appears that every third series here does not have left roots. let's filter
# these out and test against closed form.
for (parent, roots) =  filtered
  @assert roots == [closed(roots[1], n) for n in 0:24]
end

# yes, closed from matches the data. plot them all in log scale:
all = [[root, log2(parent)] for (parent, series) in filtered for root in series]
df = DataFrame(reduce(hcat, all)')
plot(df, y = :x1, Geom.line, color = :x2, Scale.y_log10, theme)

# Is the reverse tree an ordering of the natural numbers?
# -------------------------------------------------------
#
# XXX(rk): verify
#
# show that all right descendings series starting at a unique odd number
# do not have any overlapping numbers in them. this would mean that the entire
# reverse graph has unique numbers, under the condition that all left branch
# roots are also unique.
#
# try to show p2^n = q2^m where all are natural and p, q are odd and p != q:
#
# (1) p2^n = q2^m
# (2) p = q2^(m-n), if m > n then
# (3) RHS is odd * even = even, but LHS is odd. so false. if m - n = 0 then
# (4) p = q. but p != q. so false. if m - n < 0 then
# (5) p/q = 1/2^x. but then p/q would either
#     (i) have to be 1/q = 1/2^x, which is false because q is odd and 2^x is not
#     (ii) have be C * 1/2^x, but then C would have to be odd for C = q, but
#          then the denominator q would be C*2^x which is odd not even, so false
#
# therefore p2^n != q2^m.
#
# let's lint this by trying to find a brute force solution to the proof above:
#
function smash()
  p = rand(1:2:100) # p is odd
  q = rand(1:2:100) # q is odd
  m = rand(1:100)
  n = rand(1:100)

  if p == q
    return
  end

  # an optimization to only use bigint if needed. 64 signed bits, so up to 2^62.
  # can be scaled by up to 100 ~= 2^7 so let's say up to 55 is ok. use bigint
  # otherwise.
  if m > 55 || n > 55
    lhs = p * BigInt(2)^n
    rhs = q * BigInt(2)^m
  else
    lhs = p * 2^n
    rhs = q * 2^m
  end

  if lhs == rhs
    error("$p*2^$n == $q*2^$m")
  end
end

# ONE BILLIONNN, takes about 5 minutes.
for x in 1:1_000_000_000
  smash()
end

# ok. what about this: if every odd number is present in the reverse graph, then
# what would that mean?
#
# 1. all series are in the from p2^n, where p is odd. so except for when n = 0,
#    the numbers must all be even since we're multiplying them with 2's. the
#    only odd number is the root of the series. so any odd number in the reverse
#    graph must be the root of a series.
#
# 2. were interested to know if the reverse graph includes all natural numbers,
#    because then any number would be part of a tree starting with 1.
#    well in the best case, we would have all odd numbers being a root. if this
#    were true, would we then have all numbers in the tree? in other words, can
#    every natural number be expressed as p2^n where p is odd? and p and n are
#    free.
#
#    well, we've already  covered all the odd numbers, they're in.
#    we also have all the even numbers, since we can keep halving it to reach
#    either an odd number (exists), or 1, which is the 1 series.
