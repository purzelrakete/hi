using DataFrames

# Optimization
# ------------
#
# * one epoch = one forward pass and one backward pass of all the training
#   examples
# * batch size = the number of training examples in one forward/backward pass. The
#   higher the batch size, the more memory space you'll need.
# * number of iterations = number of passes, each pass using [batch size] number
#   of examples. To be clear, one pass = one forward pass + one backward pass (we
#   do not count the forward pass and backward pass as two different passes).
#
# Example: if you have 1000 training examples, and your batch size is 500, then it
# will take 2 iterations to complete 1 epoch.

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
    model.z -= opt.α * gradient(model, df)
    [i, nll(model, df)...]
  end
end

# finite difference method to validate the analytical gradients. centered version for better estimation.
finite(model::Model, df::DataFrame; ϵ::Float64 = 0.0001) = (nll(model.z + ϵ, df) - nll(model.z + ϵ, df)) / 2ϵ
