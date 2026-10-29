> ## Documentation Index
> Fetch the complete documentation index at: https://docs.z.ai/llms.txt
> Use this file to discover all available pages before exploring further.

# Official Java SDK

Z.AI Java SDK is the official Java development toolkit provided by Z.AI, offering Java developers convenient and efficient AI model integration solutions.

### Core Advantages

> **Enterprise-grade**
>
> Designed for enterprise applications, supports high concurrency and high availability

> **Easy Integration**
>
> Clean API design, comprehensive documentation, quick integration into existing projects

> **Type Safety**
>
> Complete type definitions, compile-time error checking, reducing runtime errors

> **High Performance**
>
> Optimized network request handling, supports connection pooling and asynchronous calls

### Supported Features

* **💬 Chat Conversations**: Support for single-turn and multi-turn conversations, streaming and non-streaming responses
* **🔧 Function Calling**: Enable AI models to call your custom functions
* **👁️ Vision Understanding**: Image analysis, visual understanding
* **🎨 Image Generation**: Generate high-quality images from text descriptions
* **🎬 Video Generation**: Creative content generation from text to video
* **🔊 Speech Processing**: Speech-to-text, text-to-speech
* **📊 Text Embedding**: Text vectorization, supporting semantic search
* **🤖 Intelligent Assistants**: Build professional AI assistant applications

## Technical Specifications

### Environment Requirements

* **Java Version**: Java 1.8 or higher
* **Build Tools**: Maven 3.6+ or Gradle 6.0+
* **Network Requirements**: HTTPS connection support
* **API Key**: Valid Z.AI API key required

### Dependency Management

The SDK adopts a modular design, allowing you to selectively import functional modules as needed:

* **Core Module**: Basic API calling functionality
* **Async Module**: Asynchronous and concurrent processing support
* **Utility Module**: Utility tools and auxiliary functions

## Quick Start

### Environment Requirements

> **Java Version**
>
> Java 1.8 or higher

> **Build Tools**
>
> Maven 3.6+ or Gradle 6.0+

> [!WARNING]
>
> Supports Java 8, 11, 17, 21 versions, cross-platform compatible with Windows, macOS, Linux

### Add Dependencies

**Maven**

```xml theme={null}
<dependency>
    <groupId>ai.z.openapi</groupId>
    <artifactId>zai-sdk</artifactId>
    <version>0.3.5</version>
</dependency>
```

**Gradle**

```gradle theme={null}
implementation 'ai.z.openapi:zai-sdk:0.3.5'
```

### Get API Key

1. Access [Z.AI Open Platform](https://z.ai/model-api), Register or Login.
2. Create a Standard API Key in the [API Keys](https://z.ai/manage-apikey/apikey-list) management page.
3. Copy your API Key for use.

> [!TIP]
>
> It is recommended to set the API Key as an environment variable: `export ZAI_API_KEY=your-api-key`

> [!NOTE]
>
> Domestic Z.AI platform uses ZaiClient
>
> Domestic API URL: [https://api.z.ai/api/paas/v4/](https://api.z.ai/api/paas/v4/)

#### Create Client

**Environment Variable**

```java theme={null}
import ai.z.openapi.ZaiClient;

public class QuickStart {
    public static void main(String[] args) {
        // Read API Key from environment variable
        ZaiClient client = ZaiClient.builder().ofZAI()
            .apiKey(System.getenv("ZAI_API_KEY"))
            .build();
        
        // Or use directly (if environment variable is set)
        ZaiClient client2 = ZaiClient.builder().ofZAI().build();
    }
}
```

**Direct Setting**

```java theme={null}
import ai.z.openapi.ZaiClient;

public class QuickStart {
    public static void main(String[] args) {
        // Set API Key directly
        ZaiClient client = ZaiClient.builder().ofZAI()
            .apiKey("YOUR_API_KEY")
            .build();
    }
}
```

#### Basic Conversation

```java theme={null}
import ai.z.openapi.ZaiClient;
import ai.z.openapi.service.model.*;
import ai.z.openapi.core.Constants;
import java.util.Arrays;

public class BasicChat {
    public static void main(String[] args) {
        // Initialize client
        ZaiClient client = ZaiClient.builder().ofZAI()
            .apiKey("YOUR_API_KEY")
            .build();
        
        // Create chat completion request
        ChatCompletionCreateParams request = ChatCompletionCreateParams.builder()
            .model("glm-5.3")
            .messages(Arrays.asList(
                ChatMessage.builder()
                    .role(ChatMessageRole.USER.value())
                    .content("Hello, please introduce yourself")
                    .build()
            ))
            .build();
        
        // Send request
        ChatCompletionResponse response = client.chat().createChatCompletion(request);
        
        // Get reply
        if (response.isSuccess()) {
            Object reply = response.getData().getChoices().get(0).getMessage().getContent();
            System.out.println("AI Reply: " + reply);
        } else {
            System.err.println("Error: " + response.getMsg());
        }
    }
}
```

#### Streaming Conversation

```java theme={null}
import ai.z.openapi.ZaiClient;
import ai.z.openapi.service.model.*;
import ai.z.openapi.core.Constants;
import java.util.Arrays;

public class StreamingChat {
    public static void main(String[] args) {
        ZaiClient client = ZaiClient.builder().ofZAI()
            .apiKey("YOUR_API_KEY")
            .build();
        
        // Create streaming chat request
        ChatCompletionCreateParams request = ChatCompletionCreateParams.builder()
            .model("glm-5.3")
            .messages(Arrays.asList(
                ChatMessage.builder()
                    .role(ChatMessageRole.USER.value())
                    .content("Write a poem about spring")
                    .build()
            ))
            .stream(true)
            .build();
        
        // Handle streaming response
        ChatCompletionResponse response = client.chat().createChatCompletion(request);
        
        if (response.isSuccess() && response.getFlowable() != null) {
            response.getFlowable().subscribe(
                data -> {
                    // Handle streaming data chunks
                    if (data.getChoices() != null && !data.getChoices().isEmpty()) {
                        Delta content = data.getChoices().get(0).getDelta();
                        if (content != null) {
                            System.out.print(content);
                        }
                    }
                },
                error -> System.err.println("\nStreaming error: " + error.getMessage()),
                () -> System.out.println("\nStreaming completed")
            );
        }
    }
}
```

### Complete Example

```java theme={null}
import ai.z.openapi.ZaiClient;
import ai.z.openapi.service.model.*;
import ai.z.openapi.core.Constants;
import java.util.*;

public class ChatBot {
    private final ZaiClient client;
    private final List<ChatMessage> conversation;
    
    public ChatBot(String apiKey) {
        this.client = ZaiClient.builder().ofZAI()
            .apiKey(apiKey)
            .build();
        this.conversation = new ArrayList<>();
        // Add system message
        this.conversation.add(ChatMessage.builder()
            .role(ChatMessageRole.SYSTEM.value())
            .content("You are a friendly AI assistant")
            .build());
    }
    
    public Object chat(String userInput) {
        try {
            // Add user message
        conversation.add(ChatMessage.builder()
            .role(ChatMessageRole.USER.value())
            .content(userInput)
            .build());
        
        // Create request
        ChatCompletionCreateParams request = ChatCompletionCreateParams.builder()
            .model("glm-5.3")
            .messages(conversation)
            .temperature(0.6f)
            .maxTokens(1000)
            .build();
        
        // Send request
        ChatCompletionResponse response = client.chat().createChatCompletion(request);
        
        if (response.isSuccess()) {
            // Get AI response
            Object aiResponse = response.getData().getChoices().get(0).getMessage().getContent();
            
            // Add AI response to conversation history
            conversation.add(ChatMessage.builder()
                .role(ChatMessageRole.ASSISTANT.value())
                .content(aiResponse)
                .build());
            
            return aiResponse;
        } else {
            return "Error occurred: " + response.getMsg();
        }
        
    } catch (Exception e) {
        return "Error occurred: " + e.getMessage();
        }
    }
    
    public static void main(String[] args) {
        ChatBot bot = new ChatBot(System.getenv("ZAI_API_KEY"));
        Scanner scanner = new Scanner(System.in);
        
        System.out.println("Welcome to Z.ai chatbot! Type 'quit' to exit.");
        
        while (true) {
            System.out.print("You: ");
            String input = scanner.nextLine();
            
            if ("quit".equalsIgnoreCase(input)) {
                break;
            }
            
            Object response = bot.chat(input);
            System.out.println("AI: " + response);
        }
        
        System.out.println("Goodbye!");
        scanner.close();
    }
}
```

## Advanced Features

### Function Calling

Function calling allows AI models to call functions you define to get real-time information or perform specific operations.

#### Defining and Using Functions

```java theme={null}
import ai.z.openapi.ZaiClient;
import ai.z.openapi.service.model.*;
import ai.z.openapi.core.Constants;
import java.util.*;

public class FunctionCallingExample {
    
    // Simulate weather API
    public static Map<String, Object> getWeather(String location, String date) {
        Map<String, Object> weather = new HashMap<>();
        weather.put("location", location);
        weather.put("date", date != null ? date : "today");
        weather.put("weather", "sunny");
        weather.put("temperature", "25°C");
        weather.put("humidity", "60%");
        return weather;
    }
    
    // Simulate stock API
    public static Map<String, Object> getStockPrice(String symbol) {
        Map<String, Object> stock = new HashMap<>();
        stock.put("symbol", symbol);
        stock.put("price", 150.25);
        stock.put("change", "+2.5%");
        return stock;
    }
    
    public static void main(String[] args) {
        ZaiClient client = ZaiClient.builder().ofZAI()
            .apiKey(System.getenv("ZAI_API_KEY"))
            .build();
        
        // Define function tools
        Map<String, ChatFunctionParameterProperty> properties = new HashMap<>();
        ChatFunctionParameterProperty locationProperty = ChatFunctionParameterProperty
                .builder().type("string").description("City name, for example: Beijing").build();
        properties.put("location", locationProperty);
        ChatFunctionParameterProperty unitProperty = ChatFunctionParameterProperty
                .builder().type("string").enums(Arrays.asList("celsius", "fahrenheit")).build();
        properties.put("unit", unitProperty);
        ChatTool weatherTool = ChatTool.builder()
                .type(ChatToolType.FUNCTION.value())
                .function(ChatFunction.builder()
                        .name("get_weather")
                        .description("Get weather information for a specified location")
                        .parameters(ChatFunctionParameters.builder()
                                .type("object")
                                .properties(properties)
                                .required(Collections.singletonList("location"))
                                .build())
                        .build())
                .build();
        
        // Create request
        ChatCompletionCreateParams request = ChatCompletionCreateParams.builder()
                .model("glm-5.3")
                .messages(Collections.singletonList(
                        ChatMessage.builder()
                                .role(ChatMessageRole.USER.value())
                                .content("How's the weather in Beijing today?")
                                .build()
                ))
                .tools(Collections.singletonList(weatherTool))
                .toolChoice("auto")
                .build();

        // Send request
        ChatCompletionResponse response = client.chat().createChatCompletion(request);

        if (response.isSuccess()) {
            // Handle function calling
            ChatMessage assistantMessage = response.getData().getChoices().get(0).getMessage();
            if (assistantMessage.getToolCalls() != null && !assistantMessage.getToolCalls().isEmpty()) {
                for (ToolCalls toolCall : assistantMessage.getToolCalls()) {
                    String functionName = toolCall.getFunction().getName();

                    if ("get_weather".equals(functionName)) {
                        Map<String, Object> result = getWeather("Beijing", null);
                        System.out.println("Weather info: " + result);
                    }
                }
            } else {
                System.out.println(assistantMessage.getContent());
            }
        } else {
            System.err.println("Error: " + response.getMsg());
        }
    }
}
```

## Getting Help

> **[GitHub Repository](https://github.com/zai-org/z-ai-sdk-java)**
>
> View source code, submit issues, contribute

> **[API Reference](https://docs.z.ai/api-reference)**
>
> View complete API documentation

> **[Example Projects](https://github.com/zai-org/z-ai-sdk-java/tree/main/samples)**
>
> Browse more practical application examples

> **[Best Practices](https://github.com/zai-org/z-ai-sdk-java)**
>
> Learn best practices for SDK usage

> [!NOTE]
>
> This SDK is developed based on the latest API specifications from Z.AI, ensuring synchronization with platform features. It is recommended to regularly update to the latest version for the best experience.


This documentation is built and hosted on [Mintlify](https://mintlify.com), a developer documentation platform.