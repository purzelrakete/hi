using DataFrames

# bound values in an array
function bound(a::Vector{Float64}; min = 0.0, max = 256.0)
  underflow = find(x -> x < min, a)
  overflow = find(x -> x > max, a)
  a[underflow] = min
  a[overflow] = max
  a
end

# sample rows from a dataframe
sample_df(df::DataFrame, size::Int = 1) = df[sample(1:nrow(df), size), :]

# multi dimensional sigmoid. ignores numerical issues. see scratch notebook.
sigmoid(x) = 1 ./ (1 + exp(-x))

# zscore the row vectors along the given dimension. 0 std elements have a zscore of 0.0.
znormalize(X::Matrix{Float64}; dim::Int = 1) = fillnan(0.0, zscore(X, dim))

# fill NaN values with given val
fillnan(val, X) = map(x -> isnan(x) ? val : x, X)
