using DataFrames

# x'W
type LinearNoBias <: Model
  weights::Matrix{Float64}
  LinearNoBias() = new(Matrix())
end

# optimize objective
function train(::Type{LinearNoBias}, df::DataFrame)
  model = LinearNoBias()
  image_means = by(df, :label, sdf -> mean(normalize(sdf[:image])))
  model.weights = reshape(image_means[:x1], 784, 10)
  model
end

# feed forward
function prediction(model::LinearNoBias, df::DataFrame)::DataFrame
  predictions = indmax.([normalize(x)' * model.weights for x in df[:image]])
  [df DataFrame(prediction = predictions - 1)]
end
