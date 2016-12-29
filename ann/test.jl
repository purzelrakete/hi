using Base.Test
using Distributions

include("ann.jl")

df = dataset()
@test size(df) == (70_000, 2)

# images
@test typeof(to_image(df[:image][1])) <: Image
@test typeof(to_png(df[:image][1])) == Vector{UInt8}

# image composition
svg, ctx = grid([to_png(img) for img in df[:image][1:10]])
@test typeof(svg) == Compose.SVG
@test typeof(ctx) == Compose.Context

# learning
model = train(LinearNoBias(), sample_df(df, 1000))
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

# utils
@test bound([-1, 257, 12]) == [0, 256, 12]
