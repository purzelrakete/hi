using DataFrames
using Distributions

# random k folds
immutable RandomFolds
  folds::Vector{Vector{Int}}
  RandomFolds() = new(Vector{Int}[])
end

# returns random k-folds over the dataframe.
function cv(::Type{RandomFolds}, df::DataFrame, nfolds::Int = 10)
  nrows, _ = size(df)
  sizes = foldsizing(nrows, nfolds)
  folds = RandomFolds()

  population = [1:nrows;]
  for n in sizes
    idxs = sample(1:length(population), n; replace = false)
    append!(folds.folds, [population[idxs]])
    deleteat!(population, sort(idxs))
  end

  folds
end

# returns a vector of fold sizes that sum to nrows
function foldsizing(nrows::Int, nfolds::Int)::Vector{Int}
  fold_size = floor(nrows / nfolds)
  redistribute = nrows % nfolds
  folds = ones(Integer, nfolds) * fold_size
  folds[1:redistribute] += 1
  folds
end
