using Base.Test
using Distributions

include("ann.jl")

# load the whole dataset once.
df_all = dataset()

# Dataset
# -------
#
@testset "Dataset" begin
  @test size(df_all) == (70_000, 3)
  @test ndims(df_all) == 28^2
  @test nclasses(df_all) == 10
end

# Feature transforms
# ------------------
#
@testset "Features" begin

  # binarize labels
  labels = [0, 1, 4, 1]
  df = DataFrame(y = labels)
  t = transform(BinarizeLabels, df)
  @test t.mapping == Dict(1 => 0, 2 => 1, 3 => 4)
  @test reduce(hcat, df[:y])' == [
    1 0 0;
    0 1 0;
    0 0 1;
    0 1 0]

  # reverse the transform
  untransform(t, df)
  @test df[:y] == labels

  # z normalize
  df = DataFrame(x = [[1.0, 2.0], [3.0, 4.0]])
  t = transform(ZNormalize, df)
  @test t.μs == [2.0; 3.0]
  @test round(t.σs, 1) == [1.4; 1.4]
  @test mean(reduce(hcat, df[:x])', 1)  == [0.0 0.0]
  @test std(reduce(hcat, df[:x])', 1) ≈ [1.0 1.0]

  # multiple transforms
  df = DataFrame(x = [[1.0, 1.0], [1.0, 1.0]], y = [1, 2])
  transform([BinarizeLabels, ZNormalize], df)
  df[:y] == eye(2, 2)
  df[:x] == [0.0; 0.0]
end

# Images
# -------
#
@testset "Image manipulation" begin
  @test isa(to_image(df_all[:image][1]), Image)
  @test isa(to_png(df_all[:image][1]), Vector{UInt8})

  # image composition
  svg, ctx = grid([to_png(img) for img in df_all[:image][1:10]])
  @test isa(svg, Compose.SVG)
  @test isa(ctx, Compose.Context)
end

# Training, evaluation
# --------------------
#
@testset "Training and evaluation" begin
  # learning
  model, _ = train(LinearTransform, NoopOpt(), sample_df(df_all, 1000))
  @test size(model.weights) == (28^2, 10)

  # evaluation
  df = DataFrame(y = [0, 1, 2, 1, 3], prediction = [0, 3, 2, 3, 1])
  @test confusions(df) == DataFrame(y = [1, 3], prediction = [3, 1], x1 = [2, 1])
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
end

# Cross Validation
# ----------------
#
@testset "Cross Validation" begin
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
end

# Utils
# -------
#
@testset "Utils" begin
  @test bound([-1.0, 257.0, 12.0]) == [0.0, 256.0, 12.0]
  @test mean(znormalize([1.0 1.0 0.0; -1.0 -1.0 0.0]), 1) == [0.0 0.0 0.0]
  @test std(znormalize([1.0 1.0 0.0; -1.0 -1.0 0.0]), 1) ≈ [1.0 1.0 0.0]
end

# LinearTransform
# ---------------
#
@testset "Linear Transform" begin
  df = DataFrame(x = [[1.0, 0.0], [0, 1]], y = [1; 2])
  model = LinearTransform([2 1; 2 1])

  @test ndims(model) == 2
  @test nclasses(model) == 2

  @test likelihood(model, [1.0, 0.0]) == [2.0 1.0]
  @test train(LinearTransform, NoopOpt(), df)[1].weights == eye(2, 2)
  @test prediction(model, df) == [df DataFrame(prediction = [0, 0])]

  # cvpredict
  df_sampled = sample_df(df_all, 100)
  folds = cv(RandomKFolds, df_sampled)
  Yp, models = cvpredict(LinearTransform, folds, NoopOpt(), df_sampled)

  @test length(models) == 10
  @test size(Yp) == (100, 5)
  @test by(Yp, :fold, nrow)[:x1] == ones(Int, 10) * 10
end

# BinaryLogReg
# ------------
#
@testset "Binary Logistic Regression" begin
  df = DataFrame(x = [[1.0, 0.0], [0.0, 1.0]], y = [1; 2])
  model = BinaryLogReg([1 0; 0 1])
  opt = BatchGradientDescent(0.03, 1)

  @test ndims(model) == 2
  @test nclasses(model) == 2

  @test likelihood(model, [0.0, 0.0]) == [0.5 0.5]
  @test round(nll(model, df), 2) == [1.01 0.01]
  @test size(train(BinaryLogReg, opt, df)[1].z) == (2, 2)
  @test size(train(BinaryLogReg, opt, df)[2]) == (1, 3)
  @test prediction(model, df) == [df DataFrame(prediction = [0, 1])]
  @test size(gradient(model, df)) == (2, 2)

  # cvpredict
  df_sampled = sample_df(df_all, 10)
  folds = cv(RandomKFolds, df_sampled)
  Yp, models = cvpredict(BinaryLogReg, folds, opt, df_sampled)

  @test length(models) == 10
  @test size(Yp) == (10, 5)
  @test by(Yp, :fold, nrow)[:x1] == ones(Int, 10)
end

@testset "One vs All Binary Logistic Regression" begin
  df = sample_df(df_all, 100)
  model = BinaryLogReg([1 0; 0 1])
  opt = BatchGradientDescent(0.03, 10)
  _, bt = transform([ZNormalize, BinarizeLabels], df)
  folds = cv(RandomKFolds, df)

  # cvpredict
  Yp, models, stats = cvpredict(BinaryLogReg, folds, opt, df)
  untransform(bt, Yp)

  @test length(models) == 10
  @test size(Yp) == (100, 5)
  @test by(Yp, :fold, nrow)[:x1] == ones(Int, 10) * 10
end

@testset "One vs All Binary Logistic Regression" begin
  df = sample_df(df_all, 100)
  opt = BatchGradientDescent(0.003, 10)
  _, bt = transform([ZNormalize, BinarizeLabels], df)
  folds = cv(RandomKFolds, df)

  # cvpredict
  Yp, models, stats = cvpredict(BinaryLogReg, folds, opt, df)
  untransform(bt, Yp)

  @test length(models) == 10
  @test size(Yp) == (100, 5)
  @test by(Yp, :fold, nrow)[:x1] == ones(Int, 10) * 10
end
