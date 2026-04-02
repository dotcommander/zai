The embed command generates vector embeddings using Z.AI's Embedding-3 model. Output is always JSON -- embedding vectors are numeric arrays that you pipe into downstream tools or store in a vector database.

- [Creating Embeddings](#creating-embeddings)
- [Multiple Inputs](#multiple-inputs)
- [Piping Text](#piping-text)
- [Model Override](#model-override)

<a name="creating-embeddings"></a>
## Creating Embeddings

```bash
zai embed "Hello, world!"
```

Zai returns a JSON object containing the model name, input count, token usage, and the embedding data array. Each embedding is a high-dimensional numeric vector suitable for similarity search, clustering, or classification.

<a name="multiple-inputs"></a>
## Multiple Inputs

Pass multiple strings as separate arguments to embed them in a single API request:

```bash
zai embed "first sentence" "second sentence" "third sentence"
```

Zai returns one embedding per input, each with its index in the response.

<a name="piping-text"></a>
## Piping Text

Pipe text from stdin:

```bash
echo "embed this text" | zai embed
```

You may combine stdin with positional arguments -- both are embedded together:

```bash
echo "from stdin" | zai embed "from args"
```

<a name="model-override"></a>
## Model Override

Override the model with `-m`:

```bash
zai embed "hello" -m embedding-3
```

For model configuration defaults, see the [configuration documentation](/docs/configuration#api-settings).
