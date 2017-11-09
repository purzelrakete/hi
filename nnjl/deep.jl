# Given 2x2 images, where each pixel can take a value of 0 or 1, detect
# diagonals, eg:
#
#   10
#   01

# functions
sigmoid(x) = 1 ./ (1 + exp.(-x))
relu(x) = max.(0, x)
predict(x, A, B) = sigmoid(B * relu(A * x))

# training data
n = 4
m = 2^n
X = [parse(Int, bin(x, 4)[y]) for y = 1:n, x = 0:m-1] # every 2x2 mono image
y = zeros(Int, 1, m); y[1, 6+1] = 1; y[1, 9+1] = 1

# parameter search
A = zeros(3, 4)
B = zeros(1, 3)
best_likelihood = 0.0

for i = 0:100_000
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
