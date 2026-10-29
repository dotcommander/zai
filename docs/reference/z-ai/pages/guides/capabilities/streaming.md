> ## Documentation Index
> Fetch the complete documentation index at: https://docs.z.ai/llms.txt
> Use this file to discover all available pages before exploring further.

# Streaming Messages

Streaming Messages allow real-time content retrieval while the model generates responses, without waiting for the complete response to be generated. This approach can significantly improve user experience, especially when generating long text content, as users can immediately see output beginning to appear.

## Features

Streaming messages use an incremental generation mechanism, transmitting content in chunks in real-time during the generation process, rather than waiting for the complete response to be generated before returning it all at once.

### Core Parameter Description

* **`stream=True`**: Enable streaming output, must be set to `True`
* **`model`**: Models that support streaming output, such as `glm-5.3`,  `glm-5.2`, etc.

### Response Format Description

Streaming responses use Server-Sent Events (SSE) format, with each event containing:

* `choices[0].delta.content`: Incremental text content
* `choices[0].delta.reasoning_content`: Incremental reasoning content
* `choices[0].finish_reason`: Completion reason (only appears in the last chunk)
* `usage`: Token usage statistics (only appears in the last chunk)

## Code Examples

**cURL**

```bash theme={null}
curl --location 'https://api.z.ai/api/paas/v4/chat/completions' \
--header 'Authorization: Bearer YOUR_API_KEY' \
--header 'Content-Type: application/json' \
--data '{
    "model": "glm-5.3",
    "messages": [
        {
            "role": "user",
            "content": "Write a poem about spring"
        }
    ],
    "stream": true
}'
```

**Python**

**Install SDK**

```bash theme={null}
# Install latest version
pip install zai-sdk

# Or specify version
pip install zai-sdk==0.2.3
```

**Verify Installation**

```python theme={null}
import zai
print(zai.__version__)
```

**Complete Example**

```python theme={null}
from zai import ZaiClient

# Initialize client
client = ZaiClient(api_key='Your API Key')

# Create streaming message request
response = client.chat.completions.create(
    model="glm-5.3",
    messages=[
        {"role": "user", "content": "Write a poem about spring"}
    ],
    stream=True  # Enable streaming output
)

# Process streaming response
full_content = ""
for chunk in response:
    if not chunk.choices:
        continue
    
    delta = chunk.choices[0].delta
    
    # Handle incremental content
    if hasattr(delta, 'content') and delta.content:
        full_content += delta.content
        print(delta.content, end="", flush=True)
    
    # Check if completed
    if chunk.choices[0].finish_reason:
        print(f"\n\nCompletion reason: {chunk.choices[0].finish_reason}")
        if hasattr(chunk, 'usage') and chunk.usage:
            print(f"Token usage: Input {chunk.usage.prompt_tokens}, Output {chunk.usage.completion_tokens}")

print(f"\n\nComplete content:\n{full_content}")
```

### Response Example

The streaming response format is as follows:

```
data: {"id":"1","created":1677652288,"model":"glm-5.3","choices":[{"index":0,"delta":{"content":"Spring"},"finish_reason":null}]}

data: {"id":"1","created":1677652288,"model":"glm-5.3","choices":[{"index":0,"delta":{"content":" comes"},"finish_reason":null}]}

data: {"id":"1","created":1677652288,"model":"glm-5.3","choices":[{"index":0,"delta":{"content":" with"},"finish_reason":null}]}

...

data: {"id":"1","created":1677652288,"model":"glm-5.3","choices":[{"index":0,"finish_reason":"stop","delta":{"role":"assistant","content":""}}],"usage":{"prompt_tokens":8,"completion_tokens":262,"total_tokens":270,"prompt_tokens_details":{"cached_tokens":0}}}

data: [DONE]
```

## Application Scenarios

> **Chat Applications**
>
> * Real-time conversation experience
> * Character-by-character reply display
> * Reduced waiting time

> **Content Generation**
>
> * Article writing assistant
> * Code generation tools
> * Creative content creation

> **Educational Applications**
>
> * Online Q\&A systems
> * Learning assistance tools
> * Knowledge Q\&A platforms

> **Customer Service Systems**
>
> * Intelligent customer service bots
> * Real-time problem solving
> * User support systems


This documentation is built and hosted on [Mintlify](https://mintlify.com), a developer documentation platform.