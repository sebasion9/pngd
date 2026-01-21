# pgnd

## Examples

### Draw PNG to terminal

```
go run ./examples/rgb/rgb.go /path/to/png
```
Or run in ascii mode
```
go run ./examples/rgb/rgb.go /path/to/png --ascii
```

### Draw video

Example with transcoding mp4 to avi.

```
ffmpeg -i /path/to/vid.mp4 -c:v png -an /path/output/vid.avi
go run examples/avi/avi.go /path/to/vid.avi
```

https://github.com/user-attachments/assets/fb7d232e-ed08-4006-b9ce-34f7561c2cf3


