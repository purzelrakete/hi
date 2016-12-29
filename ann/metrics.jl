using DataFrames

# one vs all evaluation
function evaluation(df; classes = [0:9; ])
  ret = DataFrame()
  ret[:, :class] = classes
  ret[:, :accuracy]  = [metric_accuracy(df, class) for class in classes]
  ret[:, :precision] = [metric_precision(df, class) for class in classes]
  ret[:, :recall]    = [metric_recall(df, class) for class in classes]
  ret
end

# metrics
metric_accuracy(df::DataFrame, class)  = round((tp(df, class) + tn(df, class)) / (1.0 * nrow(df)), 3)
metric_precision(df::DataFrame, class) = round(tp(df, class) / (1.0 * tp(df, class) + fp(df, class)), 3)
metric_recall(df::DataFrame, class)    = round(tp(df, class) / (1.0 * tp(df, class) + fn(df, class)), 3)

# binary confusions
tp(df::DataFrame, class) = countnz((df[:label] .== class) & (df[:prediction] .== class))
fn(df::DataFrame, class) = countnz((df[:label] .== class) & (df[:prediction] .!= class))
fp(df::DataFrame, class) = countnz((df[:label] .!= class) & (df[:prediction] .== class))
tn(df::DataFrame, class) = countnz((df[:label] .!= class) & (df[:prediction] .!= class))

# top confusions in a dataframe containing :label and :prediction
function confusions(df::DataFrame)
  filtered = df[df[:label] .!= df[:prediction], [:label, :prediction]]
  aggregated = by(filtered, [:label, :prediction], nrow)
  sort(aggregated, cols = order(:x1, rev = true))
end

# full confusion matrix from dataframe containing :label and :prediction
function confusion_matrix(df::DataFrame; n_classes::Int = 10)
  confusion = zeros(Int64, n_classes, n_classes)
  for row in eachrow(by(df, [:label, :prediction], nrow))
      i, j, val = row[:label] + 1, row[:prediction] + 1, row[:x1]
      confusion[i, j] = val
  end

  convert(DataFrame, confusion)
end
