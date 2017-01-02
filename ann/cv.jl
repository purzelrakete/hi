using DataFrames
using Distributions

# splits dataset into test and train partitions.
abstract Partitioner

# random k folds
immutable RandomKFolds <: Partitioner
  folds::Vector{Vector{Int}}
  n_folds::Int
end

# support collect
import Base.length
length(cv::RandomKFolds) = cv.n_folds

# iteration
Base.start(cv::RandomKFolds) = 1
Base.next(cv::RandomKFolds, k) = ([k, cvsplit(cv, k)...], k + 1)
Base.done(cv::RandomKFolds, k) = k - 1 == cv.n_folds

# train on a data partitioning scheme. do prediction on the held out sets and
# returns a single dataframe containing the fold id that each instance was in.
function cvpredict{T <: Model}(kind::Type{T}, folds::Partitioner, df::DataFrame)
  Yp = DataFrame()
  models = []
  for (k, idx_test, idx_train) in folds
      model = train(kind, df[idx_train, :])
      predictions = prediction(model, df[idx_test, :])
      predictions[:fold] = k

      Yp = [Yp; predictions]
      models = [models; model]
  end

  Yp, models
end

# split into test and train set indices for fold k
function cvsplit(cv::RandomKFolds, k::Int)
  test  = cv.folds[k]
  train = cv.folds[setdiff(1:cv.n_folds, k)]
  [test, reduce(vcat, train)]
end

# returns random k-folds over the dataframe.
function cv(::Type{RandomKFolds}, df::DataFrame, nfolds::Int = 10)
  nrows, _ = size(df)
  sizes = foldsizing(nrows, nfolds)
  folds = Vector{Int}[]

  population = [1:nrows;]
  for n in sizes
    idxs = sample(1:length(population), n; replace = false)
    append!(folds, [population[idxs]])
    deleteat!(population, sort(idxs))
  end

  RandomKFolds(folds, length(folds))
end

# returns a vector of fold sizes that sum to nrows
function foldsizing(nrows::Int, nfolds::Int)::Vector{Int}
  fold_size = floor(nrows / nfolds)
  redistribute = nrows % nfolds
  folds = ones(Integer, nfolds) * fold_size
  folds[1:redistribute] += 1
  folds
end
