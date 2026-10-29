> ## Documentation Index
> Fetch the complete documentation index at: https://docs.z.ai/llms.txt
> Use this file to discover all available pages before exploring further.

# MCP Calling

MCP server calling is based on the Model Context Protocol (MCP) standard, dynamically connecting and calling tools and resources provided by external MCP servers in `chat/completions`. Through MCP server calling, you can extend the model's capability boundaries, implementing various functions such as search, visual understanding, file processing, data analysis, and more, providing powerful tool ecosystem support for AI applications.

## Features

The MCP server calling feature integrates MCP tool calling capabilities in `chat`, supporting direct calls to `chat/completions` to connect and use various MCP servers to extend model functionality without requiring an MCP client. Through MCP server calling, AI can:

* Dynamic Tool Discovery: Automatically discover and enumerate tools provided by MCP servers
* Real-time Tool Calling: Call various tools and functions from external MCP servers
* Multi-protocol Support: Support SSE and Streamable-HTTP transport protocols
* Flexible Configuration: Support custom server addresses, authentication information, and tool permissions
* Zhipu Ecosystem: Built-in official Zhipu MCP servers, direct use without URL configuration
* Standard Compatibility: Fully compatible with MCP protocol standards, supporting third-party MCP servers

### Supported Transport Protocols

#### SSE (Server-Sent Events)

* Use Cases: Real-time data streams, long connection scenarios
* Features: Unidirectional data flow, server-initiated push
* Simple Configuration: Pass authentication information through URL parameters

#### Streamable-HTTP

* Use Cases: Standard HTTP request-response
* Features: Bidirectional communication, supports complex interactions
* Flexible Configuration: Support custom Headers and parameters

## Parameter Description

**server\_label** `string` *required*

MCP server identifier. If connecting to official Zhipu MCP servers, use `mcp code` to fill this field, no need to fill `server_url`

**server\_url** `string`

MCP server address. Required when connecting to third-party MCP servers

**transport\_type** `string` (default: streamable-http)

Transport protocol type

* `sse`: Server-Sent Events protocol
* `streamable-http`: Standard HTTP protocol (default)

**allowed\_tools** `array`

Collection of allowed tools to call. If not specified, all available tools are allowed

```json theme={null}
["webSearchPrime", "image_analysis", "video_analysis"]
```

**headers** `object`

Authentication information and custom headers required by MCP server

```json theme={null}
{
  "Authorization": "Bearer your_api_key",
  "Custom-Header": "custom_value"
}
```

## Code Examples

**cURL**

**Connect to Official Zhipu MCP Server**

```bash theme={null}
curl --location 'https://api.z.ai/api/paas/v4/chat/completions' \
--header 'Authorization: Bearer YOUR_API_KEY' \
--header 'Content-Type: application/json' \
--data '{
    "model": "glm-4.6",
    "messages": [
        {
            "role": "user",
            "content": "Search for the latest artificial intelligence development trends"
        }
    ],
    "tools": [
        {
            "type": "mcp",
            "mcp": {
                "server_label": "mcp code",
                "transport_type": "streamable-http",
                "allowed_tools": ["webSearchPrime"]
            }
        }
    ],
    "tool_choice": "auto"
}'
```

**Connect to Third-party MCP Server**

```bash theme={null}
curl --location 'https://api.z.ai/api/paas/v4/chat/completions' \
--header 'Authorization: Bearer YOUR_API_KEY' \
--header 'Content-Type: application/json' \
--data '{
    "model": "glm-4.6",
    "messages": [
        {
            "role": "user",
            "content": "Analyze the content of this document"
        }
    ],
    "tools": [
        {
            "type": "mcp",
            "mcp": {
                "server_label": "document-analyzer",
                "server_url": "https://doc-mcp.example.com/mcp",
                "transport_type": "streamable-http",
                "headers": {
                    "Authorization": "Bearer doc_analysis_token",
                    "X-Custom-Header": "custom_value"
                },
                "allowed_tools": ["analyze_document", "extract_text"]
            }
        }
    ],
    "tool_choice": "auto"
}'
```

**Using SSE Protocol**

```bash theme={null}
curl --location 'https://api.z.ai/api/paas/v4/chat/completions' \
--header 'Authorization: Bearer YOUR_API_KEY' \
--header 'Content-Type: application/json' \
--data '{
    "model": "glm-4.6",
    "messages": [
        {
            "role": "user",
            "content": "Get real-time weather information"
        }
    ],
    "tools": [
        {
            "type": "mcp",
            "mcp": {
                "server_label": "weather-service",
                "server_url": "https://weather-mcp.example.com/sse?token=weather_token",
                "transport_type": "sse",
                "allowed_tools": ["get_current_weather", "get_forecast"]
            }
        }
    ],
    "tool_choice": "auto"
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

**Connect to Official Zhipu MCP Server**

```python theme={null}
from zai import ZaiClient

# Initialize client
client = ZaiClient(api_key='Your API Key')

# Use official Zhipu MCP server (no need to configure server_url)
response = client.chat.completions.create(
    model="glm-4.6",
    messages=[
        {
            "role": "user",
            "content": "Help me search for the latest AI technology development trends"
        }
    ],
    tools=[
        {
            "type": "mcp",
            "mcp": {
                "server_label": "mcp code",
                "transport_type": "streamable-http",
                "allowed_tools": ["webSearchPrime"]  # Only allow search tools
            }
        }
    ],
    tool_choice="auto"
)

print("Search Results:")
print(response.choices[0].message.content)
```

**Connect to Third-party MCP Server**

```python theme={null}
from zai import ZaiClient

# Initialize client
client = ZaiClient(api_key='Your API Key')

# Connect to third-party MCP server
response = client.chat.completions.create(
    model="glm-4.6",
    messages=[
        {
            "role": "user",
            "content": "Analyze the content of this image"
        }
    ],
    tools=[
        {
            "type": "mcp",
            "mcp": {
                "server_label": "vision-server",
                "server_url": "https://your-mcp-server.com/mcp",
                "transport_type": "streamable-http",
                "headers": {
                    "Authorization": "Bearer your_mcp_token",
                    "Content-Type": "application/json"
                },
                "allowed_tools": ["image_analysis", "video_analysis"]
            }
        }
    ],
    tool_choice="auto"
)

print("Analysis Results:")
print(response.choices[0].message.content)
```

**Connect Using SSE Protocol**

```python theme={null}
from zai import ZaiClient

# Initialize client
client = ZaiClient(api_key='Your API Key')

# Connect to MCP server using SSE protocol
response = client.chat.completions.create(
    model="glm-4.6",
    messages=[
        {
            "role": "user",
            "content": "Get real-time stock price information"
        }
    ],
    tools=[
        {
            "type": "mcp",
            "mcp": {
                "server_label": "finance-server",
                "server_url": "https://finance-mcp.example.com/sse?Authorization=your_token",
                "transport_type": "sse",
                "allowed_tools": ["get_stock_price", "get_market_data"]
            }
        }
    ],
    tool_choice="auto"
)

print("Stock Price Information:")
print(response.choices[0].message.content)
```

## Application Scenarios

> **Intelligent Search**
>
> * Real-time information retrieval
> * Web resource search
> * News updates acquisition

> **Visual Analysis**
>
> * Image content recognition
> * Video analysis processing
> * Document image parsing

> **Data Processing**
>
> * File format conversion
> * Data cleaning and organization
> * Report generation and analysis

> **API Integration**
>
> * Third-party service calls
> * Enterprise system integration
> * Cloud service connections


This documentation is built and hosted on [Mintlify](https://mintlify.com), a developer documentation platform.