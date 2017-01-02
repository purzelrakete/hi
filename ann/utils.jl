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

# maths
sigmoid(x) = 1 ./ (1 + exp(-x))
