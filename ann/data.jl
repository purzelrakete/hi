using MNIST
using DataFrames

# merge test and train data ready for CV splitting.
function dataset()
  [df_train(); df_test()]
end

# mnist training data, 60_000 images.
function df_train()
  df_load(traindata()...)
end

# mnist test data, 10_000 images.
function df_test()
  df_load(testdata()...)
end

# convert mnist data to dataframe
function df_load(xs, ys)
  xst, yst = xs', ys'
  rows, _ = size(xst)

  DataFrame(
    image = [xst[x, :] for x in 1:rows],
    label = map(Integer, vec(yst)))
end
