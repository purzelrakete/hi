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
  images = [xst[x, :] for x in 1:rows]

  DataFrame(
    image = images,
    x = images,
    y = map(Integer, vec(yst)))
end

# accessors
ndims(df::DataFrame) = length(df[:x][1])
nclasses(df::DataFrame) = length(classes(df))
classes(df::DataFrame) = unique(df[:y])
