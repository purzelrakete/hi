using Base.Test

include("collatz.jl")
include("tree.jl")
include("viz.jl")

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
@test f_len(2) == 2
@test f_len(8) == length(f(8))
@test f_len(20) == length(f(20))
@test f_len(63_728_127) == 950

# list lengths with memo
memo = Dict{Number,Number}()
@test f_len(2; memo = memo) == 2
@test f_len(8, memo = memo) == length(f(8))
@test f_len(20; memo = memo) == length(f(20))
@test f_len(63_728_127; memo = memo) == 950

# left forward branches on doubling series 2, 4, 8, 16, etc. see the graphviz
# output to follow this example.
@test last(series(2, 100)) != 0
@test series_left_roots(2, 10) == [5, 21, 85, 341]
@test typeof(last(series_left_roots(2, 5))) == BigInt

# trees
@test Node(123).value == 123
@test map(x -> x.value, list(Node(123, right = Node(456)))) == [123, 456]
@test edges(Node(123, right = Node(456))) == [(123, 456)]
@test edges(Node(123, left = Node(789), right = Node(456))) == [(123, 789), (123, 456)]

# forward tree branching
@test hasleft(1) == false
@test hasleft(2) == false
@test hasleft(4) == false
@test hasleft(8) == false
@test hasleft(16) == true
@test left(16) == 5
@test right(16) == 32

# forward tree
@test forward(1, 0, 0) == Node(1)
@test forward(1, 0, 2) == Node(1, right = Node(2, right = Node(4)))

# viz. raises exception when the dotfile cannot be parsed.
png(forward(maxdepth = 15))
