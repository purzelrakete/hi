# binary tree
abstract Tree
immutable Empty <: Tree; end
immutable Node{T} <: Tree
  value::T
  left::Tree
  right::Tree
end

# builder with kwargs and default empty children.
Node{T}(value::T; left::Tree = Empty(), right::Tree = Empty()) = Node(
  value,
  left,
  right)

empty(tree::Node) = false
empty(tree::Empty) = true

# traversal
list(tree::Empty) = []
list(tree::Node) = [tree; list(tree.left); list(tree.right)]

# get a list of directed edges
edges(tree::Tree) = reduce(vcat, map(node_edges, list(tree)))
node_edges(x::Node) = [
  !empty(x.left) ? [(x.value, x.left.value)] : [];
  !empty(x.right) ? [(x.value, x.right.value)] : [] ]
