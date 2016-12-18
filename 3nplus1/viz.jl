# draw trees
function png(tree::Tree)
  f = open("tree.dot", "w")
  write(f, digraph(tree))
  close(f)
  run(`dot tree.dot -Tpng -O`)
end

function digraph(tree::Tree)::String
  lines = map(x -> "  $(x[1]) -> $(x[2])", edges(tree))
  joined = join(lines, ";\n")

  """
  digraph collatz{
    graph [bgcolor=transparent]
    node [fillcolor=white style=filled]
    edge [color=black]

  $joined
  }
  """
end
