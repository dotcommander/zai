The vision command analyzes images using Z.AI's vision models. Point it at a local file or a URL, optionally provide a question, and zai describes what it sees or answers your specific query about the image.

- [Analyzing an Image](#analyzing-an-image)
- [Custom Prompts](#custom-prompts)
- [Additional Options](#additional-options)

<a name="analyzing-an-image"></a>
## Analyzing an Image

Provide an image with the `-f` flag:

```bash
zai vision -f photo.jpg
```

Zai sends the image to the vision model with a default prompt ("What do you see in this image?") and prints the analysis to your terminal.

You may also analyze images from URLs:

```bash
zai vision -f https://example.com/chart.png
```

<a name="custom-prompts"></a>
## Custom Prompts

Pass a prompt as a positional argument or with `-p` to ask something specific about the image:

```bash
zai vision -f receipt.jpg "What is the total amount?"
zai vision -f chart.png -p "Explain the trends shown in this chart"
```

This focuses the model's analysis on exactly what you need.

<a name="additional-options"></a>
## Additional Options

```bash
zai vision -f photo.jpg -m glm-4v-flash
zai vision -f photo.jpg -t 0.1
```

| Flag | Effect |
|------|--------|
| `-m <model>` | Override the vision model (default configured in config) |
| `-t <value>` | Set temperature -- lower values produce more precise answers |

For model configuration defaults, see the [configuration documentation](/docs/configuration#api-settings).
