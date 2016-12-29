using Images
using Compose

# convert a datframe image entry to an image object
function to_image(pixels::Vector{Float64}; height::Int = 28)::Image
  convert(Image{Images.Gray{Images.U8}}, reshape(pixels, height, height)' / 256)
end

# rescale to l2 unit. make sure we have valid pixels.
function normalize_pixels(pixels::Vector{Float64})
  bound(round(normalize(pixels)))
end

# convert an image to png data. save and read, because reasons.
function to_png(image::Image)
  img = convert(Image{Images.Gray{Images.U8}}, image)
  filename = "$( tempname() ).png"
  save(filename, img)
  read(open(filename))
end

# straight from pixels
function to_png(pixels::Vector{Float64})
  to_png(to_image(pixels))
end

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
