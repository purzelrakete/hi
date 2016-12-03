using Gadfly

# plot a single sequence
series = f(2016)
plot(x = 1:length(series), y = series, Geom.line)

# sequence lengths
lengths = [f_len(x) for x in 2:10000]
plot(x = 1:length(lengths), y = lengths, Geom.point)

# sequence histogram
plot(x = lengths, Geom.histogram)
