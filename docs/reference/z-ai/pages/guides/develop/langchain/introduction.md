> ## Documentation Index
> Fetch the complete documentation index at: https://docs.z.ai/llms.txt
> Use this file to discover all available pages before exploring further.

# LangChain Integration

LangChain is a framework for developing applications powered by language models. Z.AI's integration with LangChain enables you to:

* Use LangChain's chain calling functionality
* Build intelligent agents and tool calling
* Implement complex conversation memory management

## Core Advantages

> **Framework Ecosystem**
>
> Access to LangChain's rich ecosystem and toolchain

> **Rapid Development**
>
> Quickly build complex AI applications using pre-built components

> **Modular Design**
>
> Flexibly combine different components to meet various needs

> **Community Support**
>
> Enjoy active open source community and rich resources

## Environment Requirements

> **Python Version**
>
> Python 3.8 or higher

> **LangChain Version**
>
> langchain\_community version 0.0.32 or higher

> [!WARNING]
>
> Please ensure langchain\_community version is 0.0.32 or higher for optimal compatibility and feature support.

## Install Dependencies

### Basic Installation

```bash theme={null}
# Install LangChain and related dependencies
pip install langchain langchainhub httpx_sse

# Install OpenAI compatible package
pip install langchain-openai
```

### Complete Installation

```bash theme={null}
# Install all related packages at once
pip install langchain langchain-openai langchainhub httpx_sse

# Verify installation
python -c "import langchain; print(langchain.__version__)"
```

## Quick Start

### Get API Key

1. Access [Z.AI Open Platform](https://z.ai/model-api), Register or Login.
2. Create a Standard API Key in the [API Keys](https://z.ai/manage-apikey/apikey-list) management page.
3. Copy your API Key for use.

> [!TIP]
>
> It is recommended to set the API Key as an environment variable: `export ZAI_API_KEY=your-api-key`

### Basic Configuration

```python theme={null}
import os
from langchain_openai import ChatOpenAI

# Create Z.AI LLM instance
llm = ChatOpenAI(
    temperature=0.6,
    model="glm-5.3",
    openai_api_key="your-Z.AI-api-key",
    openai_api_base="https://api.z.ai/api/paas/v4/"
)

# Or use environment variables
llm = ChatOpenAI(
    temperature=0.6,
    model="glm-5.3",
    openai_api_key=os.getenv("ZAI_API_KEY"),
    openai_api_base="https://api.z.ai/api/paas/v4/"
)
```

## Basic Usage Examples

### Simple Conversation

```python theme={null}
from langchain_openai import ChatOpenAI
from langchain.schema import HumanMessage, SystemMessage

# Create LLM instance
llm = ChatOpenAI(
    temperature=0.7,
    model="glm-5.3",
    openai_api_key="your-Z.AI-api-key",
    openai_api_base="https://api.z.ai/api/paas/v4/"
)

# Create messages
messages = [
    SystemMessage(content="You are a helpful AI assistant"),
    HumanMessage(content="Please introduce the development history of artificial intelligence")
]

# Call the model
response = llm(messages)
print(response.content)
```

### Using Prompt Templates

```python theme={null}
from langchain.prompts import ChatPromptTemplate
from langchain_openai import ChatOpenAI

# Create LLM
llm = ChatOpenAI(
    model="glm-5.3",
    openai_api_key="your-Z.AI-api-key",
    openai_api_base="https://api.z.ai/api/paas/v4/"
)

# Create prompt template
prompt = ChatPromptTemplate.from_messages([
    ("system", "You are a professional {domain} expert"),
    ("human", "Please explain the concept and applications of {topic}")
])

# Create chain
chain = prompt | llm

# Invoke chain
response = chain.invoke({
    "domain": "machine learning",
    "topic": "deep learning"
})

print(response.content)
```

### Conversation Memory Management

```python theme={null}
from langchain_openai import ChatOpenAI
from langchain.prompts import (
    ChatPromptTemplate,
    MessagesPlaceholder,
    SystemMessagePromptTemplate,
    HumanMessagePromptTemplate,
)
from langchain.chains import LLMChain
from langchain.memory import ConversationBufferMemory

# Create LLM
llm = ChatOpenAI(
    temperature=1.0,
    model="glm-5.3",
    openai_api_key="your-Z.AI-api-key",
    openai_api_base="https://api.z.ai/api/paas/v4/"
)

# Create prompt template
prompt = ChatPromptTemplate(
    messages=[
        SystemMessagePromptTemplate.from_template(
            "You are a nice chatbot having a conversation with a human."
        ),
        MessagesPlaceholder(variable_name="chat_history"),
        HumanMessagePromptTemplate.from_template("{question}")
    ]
)

# Create memory
memory = ConversationBufferMemory(memory_key="chat_history", return_messages=True)

# Create conversation chain
conversation = LLMChain(
    llm=llm,
    prompt=prompt,
    verbose=True,
    memory=memory
)

# Have conversations
response1 = conversation.invoke({"question": "tell me a joke"})
print("AI:", response1['text'])

response2 = conversation.invoke({"question": "tell me another one"})
print("AI:", response2['text'])
```

## Advanced Features

### Intelligent Agent

```python theme={null}
import os
from langchain import hub
from langchain.agents import AgentExecutor, create_react_agent
from langchain_community.tools.tavily_search import TavilySearchResults
from langchain_openai import ChatOpenAI

# Set search tool API key
os.environ["TAVILY_API_KEY"] = "your-tavily-api-key"

# Create LLM
llm = ChatOpenAI(
    model="glm-5.3",
    openai_api_key="your-Z.AI-api-key",
    openai_api_base="https://api.z.ai/api/paas/v4/"
)

# Create tools
tools = [TavilySearchResults(max_results=2)]

# Get prompt template
prompt = hub.pull("hwchase17/react")

# Create agent
agent = create_react_agent(llm, tools, prompt)
agent_executor = AgentExecutor(agent=agent, tools=tools, verbose=True)

# Execute task
result = agent_executor.invoke({"input": "what is LangChain?"})
print(result['output'])
```

### Streaming Output

```python theme={null}
from langchain_openai import ChatOpenAI
from langchain.callbacks.streaming_stdout import StreamingStdOutCallbackHandler
from langchain.schema import HumanMessage

# Create LLM with streaming output
llm = ChatOpenAI(
    model="glm-5.3",
    openai_api_key="your-Z.AI-api-key",
    openai_api_base="https://api.z.ai/api/paas/v4/",
    streaming=True,
    callbacks=[StreamingStdOutCallbackHandler()]
)

# Send message (output will be displayed in real-time streaming)
response = llm([HumanMessage(content="Write a poem about spring")])
```

## Best Practices

> **Performance Optimization**
>
> * Enable LangChain caching mechanism
> * Use batch processing to reduce API calls
> * Set reasonable max\_tokens limits
> * Use async processing for better concurrency

> **Error Handling**
>
> * Implement retry mechanisms and exponential backoff
> * Set reasonable timeout values
> * Log detailed error information
> * Provide fallback solutions

> **Memory Management**
>
> * Use ConversationBufferWindowMemory to limit history length
> * Regularly clean unnecessary conversation history
> * Monitor memory usage
> * Implement conversation summarization

> **Security**
>
> * Use environment variables to store API keys
> * Implement input validation and filtering
> * Monitor API usage and costs
> * Rotate API keys regularly

## Getting Help

> **[Z.AI API Documentation](https://docs.z.ai/api-reference)**
>
> View complete Z.AI API documentation

> **[LangChain Official Documentation](https://python.langchain.com/docs/get_started/introduction)**
>
> View LangChain official documentation and tutorials

> [!NOTE]
>
> LangChain is a rapidly evolving framework. It is recommended to update to the latest version regularly for optimal functionality and performance. Meanwhile, Z.AI will continue to optimize integration with LangChain to ensure the best compatibility and user experience.


This documentation is built and hosted on [Mintlify](https://mintlify.com), a developer documentation platform.