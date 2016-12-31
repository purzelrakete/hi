using Base.Test
using Distributions

include("ann.jl")

df_all = dataset()
@test size(df_all) == (70_000, 2)

# images
@test typeof(to_image(df_all[:image][1])) <: Image
@test typeof(to_png(df_all[:image][1])) == Vector{UInt8}

# image composition
svg, ctx = grid([to_png(img) for img in df_all[:image][1:10]])
@test typeof(svg) == Compose.SVG
@test typeof(ctx) == Compose.Context

# learning
model = train(LinearNoBias, sample_df(df_all, 1000))
@test size(model.weights) == (28^2, 10)

# evaluation
df = DataFrame(label = [0, 1, 2, 1, 3], prediction = [0, 3, 2, 3, 1])
@test confusions(df) == DataFrame(label = [1, 3], prediction = [3, 1], x1 = [2, 1])
@test confusion_matrix(df; n_classes = 4) == convert(DataFrame, [1 0 0 0; 0 0 0 2; 0 0 1 0; 0 1 0 0])
@test evaluation(df; classes = [0:3;]) == DataFrame(
  class     = [0:3;],
  accuracy  = [1.0, 0.4, 1.0, 0.4],
  precision = [1.0, 0.0, 1.0, 0.0],
  recall    = [1.0, 0.0, 1.0, 0.0])

# binary confusions
for i in 0:3
  @test tp(df, i) + tn(df, i) + fp(df, i) + fn(df, i) == 5
end

# cross validation
@test foldsizing(10, 10) == ones(Int64, 10)
@test foldsizing(3, 2) == [2, 1]
@test foldsizing(10, 4) == [3, 3, 2, 2]

# folds basic properties
df = DataFrame(x = [1:10;])
for i in 1:nrow(df)
  folds = cv(RandomKFolds, df, i).folds
  flattened = reduce(vcat, folds)
  @test length(flattened) == nrow(df)
  @test length(Set(flattened)) == length(flattened)
  @test length(folds) == i
end

# cv splitting
folds = RandomKFolds([[1, 2], [3, 4], [5, 6]], 3)
@test cvsplit(folds, 1) == [[1, 2], [3, 4, 5, 6]]
@test cvsplit(folds, 2) == [[3, 4], [1, 2, 5, 6]]
@test cvsplit(folds, 3) == [[5, 6], [1, 2, 3, 4]]
@test collect(folds) == [[i, cvsplit(folds, i)...] for i in 1:3]

# cvpredict
df_sampled = sample_df(df_all, 100)
folds = cv(RandomKFolds, df_sampled)
Yp, models = cvpredict(LinearNoBias, folds, df_sampled)
@test length(models) == 10
@test size(Yp) == (100, 4)
@test by(Yp, :fold, nrow)[:x1] == ones(Int, 10) * 10

# utils
@test bound([-1.0, 257.0, 12.0]) == [0.0, 256.0, 12.0]
