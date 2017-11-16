# Forward mode automatic differentiation with Dual Numbers

using DualNumbers

sigmoid(x) = 1 ./ (1 + exp.(-x))
sigmoid(Dual(0, 1))
sigmoid(Dual(-1, 1))
sigmoid(Dual(100, 1))

relu(x) = max.(0, x)
relu(Dual(10, 1))
relu(Dual(-10, 1))

predict(x, A, B) = sigmoid(B * relu(A * x))
A = dual.(randn(5, 4), ones(5, 4))
B = dual.(randn(1, 5), ones(1, 5))
predict([1 1 1 1]', A, B)
