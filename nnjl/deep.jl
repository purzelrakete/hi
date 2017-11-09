# functions
sigmoid(x) = 1 ./ (1 + exp.(-x))
relu(x) = max.(0, x)
predict(x, A, B) = sigmoid(B * relu(A * x))

# training data
X = [
  1 0 1 1
  0 1 0 1
  0 1 0 0
  1 0 0 0
]

y = [1 1 0 0]

# parameter search
A = zeros(3, 4)
B = zeros(1, 3)
best_likelihood = 0.0

for i = 0:10_000
  A_next = A + (0.5 - rand(3, 4)) * 0.1
  B_next = B + (0.5 - rand(1, 3)) * 0.1

  y_pred = predict(X, A_next, B_next)
  likelihoods = y .* y_pred + (1 - y) .* (1 - y_pred)
  likelihood = prod(likelihoods, 2)[1]

  if likelihood > best_likelihood
    println()
    println("Improved in iteration ", i)
    println("Likelihoods: ", likelihoods)
    println("Likelihood: ", likelihood, ". NLL: ", -log.(likelihood))
    A = A_next
    B = B_next
    best_likelihood = likelihood
  end
end
