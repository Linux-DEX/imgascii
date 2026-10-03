# imgascii

Turn a picture into plain ASCII art: characters you can select, copy, and paste into a file, chat, or editor. There are no color codes.

Dark pixels become dense characters (`$`, `@`, `#`). Light pixels become sparse ones, down to a space. A terminal cell is about twice as tall as it is wide, so the row count is half of a square pixel grid and the picture keeps its shape.

## Requirements

- Go (this module uses the version in `go.mod`)
- A monospace font when you view or paste the art

## Build

From this directory:

```bash
go build -o imgascii .
```

Install it onto your `PATH`:

```bash
go install .
```

`go install` puts the binary in `$(go env GOPATH)/bin`. That directory needs to be on your `PATH`.

Run without installing. The image path or URL comes after `.`:

```bash
go run . photo.png
go run . --width 80 ~/Pictures/logo.webp
go run . https://example.com/logo.png
go run . --width 60 https://example.com/banner.jpg
```

## Usage

```bash
imgascii [--width N] <file-or-url>
```

One argument, either a local file or an `http`/`https` URL. The art is written to stdout. Errors go to stderr.

```bash
imgascii photo.png
imgascii --width 80 ~/Pictures/logo.webp
imgascii https://example.com/logo.png
imgascii --width 60 https://example.com/banner.jpg
```

Help:

```bash
imgascii -h
```

### Width

`--width` is the number of columns. The height is chosen from the image aspect ratio, so the picture is not stretched.

When you leave `--width` off, imgascii uses the width of the terminal. That only works when stdout is a terminal. Piping or redirecting needs an explicit width:

```bash
imgascii --width 80 photo.png > art.txt
imgascii --width 40 photo.png | less
```

`art.txt` is plain text. Open it or copy it from the terminal.

`--width` must be at least 1. `0` means "use the terminal width."

### Local files

Pass a path. The format is detected from the file bytes, not the extension, so a file named `image` or `photo.bin` still works if the contents are a supported image.

### URLs

`http://` and `https://` are downloaded. Anything else is opened as a local path.

- Timeout: 15 seconds
- Largest download: 32MiB
- Redirects: up to the Go HTTP client's default (10)
- Only a `200` response is decoded

## Formats

| Format | Supported |
| --- | --- |
| JPEG | yes |
| PNG | yes |
| GIF | yes, first frame only |
| BMP | yes |
| TIFF | yes |
| WebP | yes |
| AVIF, HEIC, SVG | no |

GIF animation is not played. The first frame is the picture.

## How the picture is drawn

1. The image is scaled with Catmull-Rom to `width` characters across. The number of rows follows the source aspect ratio, then is halved because a character cell is taller than it is wide. There is always at least one row.
2. Each pixel becomes one ASCII character. Brightness uses the usual weights: 30% red, 59% green, 11% blue.
3. The ramp runs from a space (white) to `$` (black):

```text
 .'`^",:;Il!i><~+_-?][}{1)(|\/tfjrxnuvczXYUJCLQ0OZmwqpdbkhao*#MW&8%B@$
```

4. A pixel with alpha below 16 is a space.

The result is only those characters and newlines. No escape codes.

## Examples

A photo at the full terminal width:

```bash
imgascii vacation.jpg
```

A small logo, 40 columns:

```bash
imgascii --width 40 logo.png
```

Save text you can copy:

```bash
go run . --width 80 photo.png > photo.txt
go run . --width 80 https://example.com/logo.png > logo.txt
```

## Tests

```bash
go test ./...
```
