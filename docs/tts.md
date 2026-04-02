The tts command converts text to speech using Z.AI's GLM-TTS model. Pass it a string or pipe in text, and zai generates an audio file that downloads to your current directory.

- [Generating Speech](#generating-speech)
- [Choosing a Voice](#choosing-a-voice)
- [Output Options](#output-options)

<a name="generating-speech"></a>
## Generating Speech

```bash
zai tts "Hello, welcome to the future of AI."
```

Zai generates the audio and saves it to `zai-tts-{timestamp}.wav` in the current directory, printing the file path and size when complete.

You may pipe text from stdin:

```bash
echo "Welcome to the conference" | zai tts
```

<a name="choosing-a-voice"></a>
## Choosing a Voice

Select from seven available voices with `--voice`:

| Voice | Description |
|-------|-------------|
| `tongtong` | Default voice |
| `xiaochen` | Alternative voice |
| `chuichui` | Alternative voice |
| `jam` | Alternative voice |
| `kazi` | Alternative voice |
| `douji` | Alternative voice |
| `luodo` | Alternative voice |

```bash
zai tts "Good morning" --voice xiaochen
```

The default voice is configured in your config file under `tts.voice`. For details, see the [configuration documentation](/docs/configuration#tts-settings).

<a name="output-options"></a>
## Output Options

```bash
zai tts "Hello" -o greeting.wav
zai tts "Hello" --format pcm
zai tts "Hello" --speed 2
zai tts "Hello" -m glm-tts
```

| Flag | Effect |
|------|--------|
| `-o <path>` | Custom output path instead of auto-generated filename |
| `--format <fmt>` | Audio format: `wav` (default) or `pcm` |
| `--speed <n>` | Playback speed multiplier (e.g., `2` for double speed) |
| `-m <model>` | Override the TTS model |

The default format is configured under `tts.response_format` in your config. For details, see the [configuration documentation](/docs/configuration#tts-settings).
