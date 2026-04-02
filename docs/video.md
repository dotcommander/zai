The video command generates videos using Z.AI's CogVideoX-3 model. Video generation is asynchronous -- zai submits the request, polls until the result is ready (typically 1-3 minutes), and auto-downloads the file to your current directory.

- [Creating a Video](#creating-a-video)
- [Image-to-Video](#image-to-video)
- [Video Options](#video-options)
- [Available Sizes](#available-sizes)

<a name="creating-a-video"></a>
## Creating a Video

```bash
zai video "A cat playing with a ball of yarn"
```

Zai displays the task ID and polls for completion, showing elapsed time as it waits. Once the video is ready, it downloads to `zai-video-{timestamp}.mp4` in the current directory.

> [!NOTE]
> Video generation typically takes 1-3 minutes. Zai handles the polling automatically. You can extend the polling timeout with `--poll-timeout` if needed.

<a name="image-to-video"></a>
## Image-to-Video

Animate a still image by providing it with `-f`:

```bash
zai video -f photo.jpg "Make the clouds move slowly"
```

Create a transition between two frames by providing both:

```bash
zai video -f first.jpg -f last.jpg "Smooth camera pan from left to right"
```

The first file becomes the opening frame and the second becomes the closing frame. Zai generates the motion between them.

<a name="video-options"></a>
## Video Options

```bash
zai video "prompt" -q quality
zai video "prompt" -s 1280x720
zai video "prompt" --fps 60
zai video "prompt" --duration 10
zai video "prompt" --with-audio
zai video "prompt" -o my-video.mp4
zai video "prompt" --show
zai video "prompt" --poll-timeout 5m
```

| Flag | Effect |
|------|--------|
| `-q <level>` | Set quality: `speed` (default) or `quality` |
| `-s <size>` | Output resolution (default: `1920x1080`) |
| `--fps <n>` | Frames per second (default: `30`) |
| `--duration <n>` | Video length in seconds (default: `5`) |
| `--with-audio` | Generate AI sound effects for the video |
| `-o <path>` | Custom output path instead of auto-generated filename |
| `--show` | Open in your default video player after generation |
| `--poll-timeout <dur>` | Maximum time to wait for generation (default: `3m`) |

> [!NOTE]
> Pricing is approximately $0.20 per video.

<a name="available-sizes"></a>
## Available Sizes

Supported resolutions include `1280x720`, `1024x1024`, `1920x1080`, `3840x2160`, and others. The default is `1920x1080`.

For model configuration defaults, see the [configuration documentation](/docs/configuration#api-settings).
