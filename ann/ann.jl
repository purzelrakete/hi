using DataFrames

# a parametric model
abstract Model
likelihood(model::Model, df::DataFrame) = [likelihood(model, x) for x in df[:image]]
prediction(model::Model, df::DataFrame) = [df DataFrame(prediction = indmax.(likelihood(model, df)) - 1)]
nll(model::Model, df::DataFrame) = sum([nll(model, df[i, :image], df[i, :label]) for i in 1:length(df)])
gradient(model::Model, df::DataFrame, j::Int) = sum([gradient(model, x[:image], x[:label], j) for x in eachrow(df)])
gradient(model::Model, df::DataFrame, j::Int, class::Int) = sum([gradient(model, x[:image], map(signed, x[:label] .== class), j) for x in eachrow(df)])

# type for the list of all pixels in an image
typealias Pixels Vector{Float64}

# includes
include("cv.jl")
include("data.jl")
include("images.jl")
include("metrics.jl")
include("models.jl")
include("utils.jl")
