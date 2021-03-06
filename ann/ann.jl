using DataFrames

# an optimization method
abstract Optimizer

# a parametric model
abstract Model
likelihood(model::Model, df::DataFrame) = [likelihood(model, x) for x in df[:x]]
prediction(model::Model, df::DataFrame) = [df DataFrame(prediction = indmax.(likelihood(model, df)) - 1)]
nll(model::Model, df::DataFrame) = sum([nll(model, r[:x], r[:y]) for r in eachrow(df)])
gradient(model::Model, df::DataFrame) = sum([gradient(model, r[:x], r[:y]) for r in eachrow(df)])

# type for the list of all pixels in an image
typealias Pixels Vector{Float64}

# includes
include("cv.jl")
include("data.jl")
include("features.jl")
include("images.jl")
include("metrics.jl")
include("models.jl")
include("optim.jl")
include("utils.jl")
