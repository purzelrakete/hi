using DataFrames

# transform the data, with reverse. lifted from sklearn
abstract Transform

# apply a list of transforms
transform(ts::Vector{DataType}, df::DataFrame) = [transform(t, df) for t in ts]

# reverse a list of transforms
untransform(ts::Vector{Transform}, df::DataFrame) = [untransform(t, df) for t in ts]

type BinarizeLabels <: Transform
  mapping::Dict{Int,Int}
end

# one hot encode the labels in :y. stores mapping of index to class.
function transform(::Type{BinarizeLabels}, df::DataFrame)
  cls, ncls, nrows = classes(df), nclasses(df), size(df)[1]
  fwd = Dict(zip(1:ncls, cls))
  backwd = Dict(zip(cls, 1:ncls))

  binarized = zeros(Int, nrows, ncls)
  for (i, j) in enumerate(df[:y])
    binarized[i, backwd[j]] = 1
  end

  # transform
  df[:y] = [binarized[x, :] for x in 1:nrows]

  BinarizeLabels(fwd)
end

type ZNormalize <: Transform
  μs::Vector{Float64}
  σs::Vector{Float64}
end

# z-score on vectors :x. stores original feature means and standard deviations.
function transform(::Type{ZNormalize}, df::DataFrame)
  X, nrows, dim = reduce(hcat, df[:x])', size(df)[1], 1
  normalized = znormalize(X; dim = dim)

  # transform
  df[:x] = [normalized[x, :] for x in 1:nrows]

  # calculate stats on input and convert to vectors
  μs = mean(X, dim)'[:, 1]
  σs = std(X, dim)'[:, 1]

  ZNormalize(μs, σs)
end
