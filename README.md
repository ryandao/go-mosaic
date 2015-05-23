## go-mosaic

A library for creating a photo mosaic from a target photo and a number of "tile" images.

### Installation

    go get github.com/ryandao/gomosaic

### Example usage

    import "github.com/ryandao/gomosaic"
    import "image/jpeg"

    result, _ := gomosaic.Mosaic(target, seeds, sqSize)

    // Save output image to file
    out, _ := os.Create("./output.jpg")
    err = jpeg.Encode(out, output, &jpeg.Options{Quality: 80})



