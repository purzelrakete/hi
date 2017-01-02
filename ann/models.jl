using DataFrames

# x'W
immutable LinearTransform <: Model
  weights::Matrix{Float64}
end

likelihood(model::LinearTransform, x::Pixels) = normalize(x)' * model.weights

# optimize objective
function train(::Type{LinearTransform}, df::DataFrame; n_pixels = 784, n_classes = 10)::Model
  means = by(df, :label, sdf -> mean(normalize(sdf[:image])))[:x1]
  LinearTransform(reshape(means, n_pixels, n_classes))
end

# logistic regression
type BinaryLogReg <: Model
  z::Vector{Float64}
  bias::Float64
end

likelihood(model::BinaryLogReg, x::Pixels) = sigmoid(x' * model.z + model.bias)

# train a model
function train(::Type{BinaryLogReg}, df::DataFrame)::Model
  "NO"
end
