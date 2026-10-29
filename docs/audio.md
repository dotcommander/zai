The audio command transcribes speech to text using Z.AI's GLM-ASR model. Hand it an audio file, a YouTube URL, or piped audio data, and zai returns the transcription. Long files are automatically split into chunks and transcribed in parallel.

- [Transcribing a File](#transcribing-a-file)
- [Domain Vocabulary](#domain-vocabulary)
- [YouTube Transcription](#youtube-transcription)
- [Supported Formats](#supported-formats)
- [Preprocessing](#preprocessing)
- [Additional Options](#additional-options)

<a name="transcribing-a-file"></a>
## Transcribing a File

```bash
zai audio -f recording.wav
```

Zai transcribes the file and prints the text to stdout. You may also pipe audio from stdin:

```bash
cat recording.wav | zai audio
```

<a name="domain-vocabulary"></a>
## Domain Vocabulary

When transcribing technical content, provide domain-specific vocabulary with `--hotwords` so the model recognizes specialized terms:

```bash
zai audio -f lecture.wav --hotwords "kubernetes,docker,etcd,kubectl"
```

You may specify up to 100 hotwords, comma-separated. This significantly improves accuracy for jargon-heavy recordings.

<a name="youtube-transcription"></a>
## YouTube Transcription

Transcribe audio directly from YouTube videos with the `--video` flag. This requires [yt-dlp](https://github.com/yt-dlp/yt-dlp) to be installed:

```bash
zai audio --video https://youtu.be/abc123
```

Combine with Voice Activity Detection to strip silence and reduce API costs:

```bash
zai audio --video https://youtu.be/abc123 --vad
```

VAD removes silent segments before sending audio to the API, which reduces both processing time and cost.

<a name="supported-formats"></a>
## Supported Formats

| Format | Extension |
|--------|-----------|
| WAV | `.wav` |
| MP3 | `.mp3` |
| MP4 | `.mp4` |
| M4A | `.m4a` |
| FLAC | `.flac` |
| AAC | `.aac` |
| Ogg | `.ogg` |

The maximum request file size is 25 MiB. Larger files are automatically split into 25-second chunks and transcribed in parallel without streaming. Local input files are preserved; only generated or downloaded temporary files are removed.

Use `zai audio -f recording.wav --stream` to print transcription text as it arrives. Streaming cannot be combined with `--json` and requires inputs no larger than 25 MiB; rerun without `--stream` to process larger files in parallel chunks.

<a name="preprocessing"></a>
## Preprocessing

Audio preprocessing is enabled by default. Zai converts files to 16kHz mono WAV, which is the optimal format for the ASR model. This requires [ffmpeg](https://ffmpeg.org/) to be installed.

Disable preprocessing with `--preprocess=false` if your files are already in the correct format.

<a name="additional-options"></a>
## Additional Options

```bash
zai audio -f speech.mp3 --json
zai audio -f speech.mp3 -l en
zai audio -f speech.mp3 -p "Previous context"
zai audio -f speech.mp3 --resume
zai audio -f speech.mp3 --clear-cache
zai audio -f speech.mp3 -m glm-asr-2512
```

| Flag | Effect |
|------|--------|
| `--json` | Output transcription as JSON |
| `--stream` | Print text deltas as they arrive (up to 25 MiB; incompatible with `--json`) |
| `-l <lang>` | Specify the audio language |
| `-p <text>` | Provide prior context to improve accuracy |
| `--resume` | Resume a partial transcription |
| `--clear-cache` | Clear cached results and start fresh |
| `-m <model>` | Override the audio model |

For model configuration defaults, see the [configuration documentation](/docs/configuration#api-settings).
