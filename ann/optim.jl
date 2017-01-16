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
    model.z += opt.α * gradient(model, df)
    [i, nll(model, df)...]
  end
end
