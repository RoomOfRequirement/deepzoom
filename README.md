#

## deepzoom tool written in golang

this project is inspired by [deepzoom python tool from openzoom](https://github.com/openzoom/deepzoom.py/blob/develop/deepzoom.py)

## usage

```bash
Deepzoom tool in golang to generate dzi files
Version: 0.0.1

Usage: deepzoom [-hvsftld] [-h help] [-v version] [-s source image file path] [-f source image format(jpg or png)] [-t tile size] [-l tile overlap] [-q output image quality] [-d destination of dzi file path]

Options
  -d string
        destination of dzi file path
  -f string
        source image format, it should be jpg or png (default "jpg")
  -h    help info
  -l int
        tile overlap
  -q float
        output image quality (default 0.8)
  -s string
        source image file path
  -t int
        tile size (default 256)
  -v    version info

```