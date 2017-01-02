using Images
using Compose

# a single mnist image
typealias MnistImage Image{Images.Gray{Images.U8}}

# convert a dataframe image entry to an image object
function to_image(pixels::Pixels; height::Int = 28)::Image
  convert(MnistImage, reshape(pixels, height, height)' / 256)
end

# convert an image to png data. save and read, because reasons.
function to_png(image::Image)
  img = convert(MnistImage, image)
  filename = "$( tempname() ).png"
  save(filename, img)
  read(open(filename))
end

# rescale to l2 unit. make sure we have valid pixels.
normalize_pixels(pixels::Pixels) = bound(round(normalize(pixels)))
to_png(pixels::Pixels) = to_png(to_image(pixels))

# compose a bunch of pngs images into a grid
function grid(images; n_cols::Integer = 30, img_width = 13mm)
  bitmaps = []
  n_rows = Integer(ceil(length(images) / n_cols))
  for i = 1:n_rows
      for j = 1:n_cols
          index = (i - 1) * n_cols + j
          if index <= length(images)
              left = img_width * (j - 1)
              right = img_width * i
              width = img_width + 1mm
              height = img_width + 1mm
              png = bitmap("image/png", images[index], left, right, width, height)
              append!(bitmaps, [png])
          end
      end
  end

  width = n_cols * img_width
  height = n_rows * img_width + 2cm
  SVG(width, height), compose(context(), bitmaps...)
end
