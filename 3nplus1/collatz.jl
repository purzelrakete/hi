# next number in the series
function succ(n::Int64)::Int64
  if n < 2
    error("only natural numbers > 1 are allowed.")
  elseif n % 2 == 0
    n / 2
  else
    n * 3 + 1
  end
end

# returns the series starting with n and ending with 1
f(n) = f(n, Int64[])
function f(n::Int64, previous)
  if n == 1
    [previous; 1]
  else
    f(succ(n), [previous; n])
  end
end

# number of steps to 1. exploit FACTS.
function f_len(n)::Int64
  i = 1

  while 1 != if n % 2 == 0 # even
      i += 1
      n = n / 2
    else # odd number * 3 is odd (odd * odd). odd + 1 -> even. take 2 steps:
      i += 2
      n = (n * 3 + 1) / 2
    end
  end

  i
end

# forward
left(n::Number)::Integer = (n - 1) / 3
right(n::Number) = n * 2
hasleft(n::Number)::Bool = (n - 1) % 3 == 0 && # be integral
  left(n) % 2 == 1 && # be odd
  left(n) > 1 # be positive, not 1

# looking at the series formed by right branches starting from a node q,
# we will find it to be in the form q*2^i.
series(q::Number, depth::Int64 = 10) = [q * BigInt(2)^i for i in 0:depth]

# following a series, what are the left branches that are immediate
# children of nodes on the double series?
series_left_roots(root::Number, depth::Int64 = 10) = [
  left(x) for x in series(root, depth) if hasleft(x) ]

# returns the forward tree
forward(; maxdepth = 10) = forward(1, 0, maxdepth)
function forward(root::Int64, depth::Int64, maxdepth::Int64)::Tree
  if depth > maxdepth
    Empty()
  else
    Node(
      root,
      hasleft(root) ? forward(left(root), depth + 1, maxdepth) : Empty(),
      forward(right(root), depth + 1, maxdepth)
    )
  end
end
