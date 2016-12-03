using Base.Test

# succ
@test succ(2) == 1
@test succ(10) == 5
@test succ(11) == 34

# succ bounds
@test_throws ErrorException succ(-1)
@test_throws ErrorException succ(0)
@test_throws ErrorException succ(1)

# lists
@test f(8) == [8, 4, 2, 1]
@test f(20) == [20, 10, 5, 16, 8, 4, 2, 1]

# list length
@test f_len(8) == length(f(8))
@test f_len(20) == length(f(20))
