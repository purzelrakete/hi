# Given 2x2 images, where each pixel can take a value of 0 or 1, detect
# diagonals, eg:
#
#   1 0
#   0 1

using AutoGrad
using DataFrames
using Gadfly

sigmoid(x) = 1 ./ (1 + exp.(-x))
relu(x) = max.(0, x)
predict(X, θ) = sigmoid(θ[3] * relu(θ[1] * X + θ[2]) + θ[4])

function loss(X, y, θ)
  y_pred = predict(X, θ)
  nll = -log.(y .* y_pred + (1 - y) .* (1 - y_pred))
  sum(nll, 2)[1]
end

# training data
n = 2^2
m = 2^n
X = [parse(Int, bin(x, 4)[y]) for y = 1:n, x = 0:m-1] # every 2x2 mono image
y = zeros(Int, 1, m); [y[1, x+1] = 1 for x in [6, 9]] # diagonals in 6 and 9

# hyperparameters
α = 0.05

# parameters
θ = [
  randn(5, 4) * sqrt(2 / 4), 0,
  randn(1, 5) * sqrt(2 / 5), 0 ]

losses = []
grads = grad(loss, 3)
for i = 0:10_000
  l = loss(X, y, θ)
  g = grads(X, y, θ)

  for p in 1:length(θ)
    θ[p] = θ[p] - α * g[p]
  end

  if mod(i, 100) == 0
    println("Loss: ", @sprintf("%.3f", l), " Likelihood: ", @sprintf("%.3f", exp(-l)))
    push!(losses, l)
  end
end

df = DataFrame(likelihood = exp.(-losses))
plot(df, y = :likelihood, Geom.line)
