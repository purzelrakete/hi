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

# number of steps to 1
function f_len(n)::Int64
  i = 1
  while succ(n) != 1
    n = succ(n)
    i += 1
  end

  i + 1
end
