Zai generates images using Z.AI's CogView-4 model. Give it a short description and it enhances your prompt with professional lighting, composition, and style directives before generating the image. The result auto-downloads to your current directory.

- [Creating an Image](#creating-an-image)
- [AI Prompt Enhancement](#ai-prompt-enhancement)
- [Disabling Enhancement](#disabling-enhancement)
- [Image Sizes](#image-sizes)
- [Additional Options](#additional-options)

<a name="creating-an-image"></a>
## Creating an Image

```bash
zai image "a wizard casting lightning in a crystal tower"
```

Zai enhances your prompt, sends it to CogView-4, and downloads the result to `zai-image-{timestamp}.png` in the current directory. The terminal output shows both the original and enhanced prompts, the image URL, and the local file path.

<a name="ai-prompt-enhancement"></a>
## AI Prompt Enhancement

By default, zai transforms your prompt before sending it to the image model. A short description like "a wizard" becomes a detailed, cinematic prompt with lighting, composition, and style directives. The original and enhanced prompts are combined for best results.

You can see the enhanced prompt in the terminal output whenever enhancement is active.

<a name="disabling-enhancement"></a>
## Disabling Enhancement

If you want exact control over the prompt, disable enhancement:

```bash
zai image "minimalist black circle on white background" --no-enhance
```

This sends your prompt verbatim to CogView-4.

<a name="image-sizes"></a>
## Image Sizes

Set the output size with `-s`:

```bash
zai image "landscape" -s 1024x1024
zai image "portrait" -s 768x1344
zai image "panorama" -s 1344x768
zai image "banner" -s 1024x768
```

The default is `1024x1024` (square).

<a name="additional-options"></a>
## Additional Options

```bash
zai image "sunset" -o sunset.png
zai image "sunset" --show
zai image "sunset" --copy
zai image "sunset" -q standard
zai image "sunset" -m cogview-4-250304
```

| Flag | Effect |
|------|--------|
| `-o <path>` | Custom output path instead of auto-generated filename |
| `--show` | Open the image in your default viewer after generation |
| `--copy` | Copy the image URL to the clipboard |
| `-q <quality>` | Set quality to `standard` or `hd` (default: `hd`) |
| `-m <model>` | Override the image model |

For model configuration defaults, see the [configuration documentation](/docs/configuration#api-settings).
