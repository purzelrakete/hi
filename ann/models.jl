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
  z::Matrix{Float64}
end

likelihood(model::BinaryLogReg, x::Pixels) = sigmoid(x' * model.z)
nll(model::BinaryLogReg, x::Pixels, y::Int) = -y * (x' * model.z) + log(1 + exp(x' * model.z))
gradient(model::BinaryLogReg, x::Pixels, y::Int, j::Int) = ((x[j] * exp(x' * model.z)) ./ (1 + exp(x' * model.z))) - y * x[j]

# train a model
train(::Type{BinaryLogReg}, df::DataFrame)::Model = error("not implemented")
