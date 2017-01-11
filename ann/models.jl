using DataFrames

# Linear Transform
# ----------------
#
# x'W
#
immutable LinearTransform <: Model
  weights::Matrix{Float64}
end

# attributes
ndims(model::LinearTransform) = size(model.weights)[1]
nclasses(model::LinearTransform) = size(model.weights)[2]

# fitting
likelihood(model::LinearTransform, x::Pixels) = normalize(x)' * model.weights

# fit
function train(::Type{LinearTransform}, opt::Optimizer, df::DataFrame)
  means = by(df, :y, sdf -> mean(normalize(sdf[:x])))[:x1]
  model = LinearTransform(reshape(means, ndims(df), nclasses(df)))
  model, DataFrame() # return empty stats
end

# Logistic Regression
# -------------------
#
# sigmoid(x'W)
#
type BinaryLogReg <: Model
  z::Matrix{Float64}
end

# attributes
ndims(model::BinaryLogReg) = size(model.z)[1]
nclasses(model::BinaryLogReg) = size(model.z)[2]

# fitting
likelihood(model::BinaryLogReg, x::Pixels) = sigmoid(x' * model.z)
nll(model::BinaryLogReg, x::Pixels, y::Int) = -y * (x' * model.z) + log(1 + exp(x' * model.z))
gradient(model::BinaryLogReg, x::Pixels, y::Int, j::Int) = x[j] * sigmoid(x' * model.z) - (y * x[j])

# fit
function train(::Type{BinaryLogReg}, opt::Optimizer, df::DataFrame)
  model = BinaryLogReg(rand(ndims(df), nclasses(df)))
  stats = optimize(opt, model, df)
  model, convert(DataFrame, reduce(hcat, stats)')
end
