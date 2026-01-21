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

<p align="center" width="100%">
<video src="https://git.franzlla.ng/sebasion/pngd/assets/badapple.mp4" width="80%" controls></video>
</p>

