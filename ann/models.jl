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
  model, DataFrame() # return no stats. didn't learn anything.
end

# Logistic Regression
# -------------------
#
# sigmoid(x'W)
#
type BinaryLogReg <: Model
  z::Matrix{Float64} # double precision
end

# attributes
ndims(model::BinaryLogReg) = size(model.z)[1]
nclasses(model::BinaryLogReg) = size(model.z)[2]

# fitting. works for {0, 1} y's as well as one hot encoded y's.
likelihood(model::BinaryLogReg, x::Pixels) = sigmoid(x' * model.z)
nll(model::BinaryLogReg, x::Pixels, y::Union{Int,Array{Int}}) = log(1 + exp(x' * model.z)) - y' .* (x' * model.z)
gradient(model::BinaryLogReg, x::Pixels, y::Union{Int,Array{Int}}) = x * (sigmoid(x' * model.z) - y')

# fit
function train(::Type{BinaryLogReg}, opt::Optimizer, df::DataFrame)
  model = BinaryLogReg(rand(ndims(df), nclasses(df)))
  stats = @time optimize(opt, model, df)
  model, convert(DataFrame, reduce(hcat, stats)')
end
