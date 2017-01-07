using DataFrames

# For basic models such as LinearTransform
type NoopOpt <: Optimizer; end

# gradient descent over all training examples
type BatchGradientDescent <: Optimizer
  α::Float64
  max_iterations::Int
end

# full gradient descent, until max_iterations has been reached.
function optimize(opt::BatchGradientDescent, model::Model, df::DataFrame)
  map(1:opt.max_iterations) do i
    gradient_update(model, df, opt.α)
    [i, nll(model, df)...]
  end
end

# update the current best solution by a single step of gradient descent.
function gradient_update(model::Model, df::DataFrame, α::Float64 = 0.03)
  for j in 1:ndims(model)
    δ = α * gradient(model, df, j, 0)

    # FIXME(rk): this only works with the lr model.
    model.z[j, :] += δ'
  end
end
