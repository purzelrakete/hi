using MNIST
using DataFrames

# load mnist
dataset() = [df_train(); df_test()]
df_train() = df_load(traindata()...)
df_test() = df_load(testdata()...)

# convert mnist data to dataframe
function df_load(xs, ys)
  xst, yst = xs', ys'
  rows, _ = size(xst)

  DataFrame(
    image = [xst[x, :] for x in 1:rows],
    label = map(Integer, vec(yst)))
end

ndims(df::DataFrame) = length(df[:image][1])
nclasses(df::DataFrame) = length(classes(df))
classes(df::DataFrame) = unique(df[:label])
