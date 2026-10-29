> ## Documentation Index
> Fetch the complete documentation index at: https://docs.z.ai/llms.txt
> Use this file to discover all available pages before exploring further.

# Tokenizer

> `Tokenizer` is used to split text into `tokens` recognizable by the model and calculate the count. It receives user input text, processes it through the model for tokenization, and finally returns the corresponding `token` count. It is suitable for text length evaluation, model input estimation, dialogue context truncation, cost calculation, etc.



## OpenAPI

````yaml POST /paas/v4/tokenizer
openapi: 3.0.1
info:
  title: Z.AI API
  description: Z.AI API available endpoints
  license:
    name: Z.AI Developer Agreement and Policy
    url: https://chat.z.ai/legal-agreement/terms-of-service
  version: 1.0.0
  contact:
    name: Z.AI Developers
    url: https://chat.z.ai/legal-agreement/privacy-policy
    email: user_feedback@z.ai
servers:
  - url: https://api.z.ai/api
    description: Production server
security:
  - bearerAuth: []
paths:
  /paas/v4/tokenizer:
    post:
      summary: Text Tokenizer
      description: >-
        `Tokenizer` is used to split text into `tokens` recognizable by the
        model and calculate the count. It receives user input text, processes it
        through the model for tokenization, and finally returns the
        corresponding `token` count. It is suitable for text length evaluation,
        model input estimation, dialogue context truncation, cost calculation,
        etc.
      requestBody:
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/TokenizerRequest'
            examples:
              Text Tokenization Example:
                value:
                  model: glm-4.6
                  messages:
                    - role: user
                      content: >-
                        What opportunities and challenges will the Chinese large
                        model industry face in 2025?
        required: true
      responses:
        '200':
          description: Business processing successful
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/TokenizerResponse'
        default:
          description: The request has failed.
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
components:
  schemas:
    TokenizerRequest:
      type: object
      required:
        - model
        - messages
      properties:
        model:
          type: string
          description: The model code to be called.
          example: glm-4.6
          default: glm-4.6
          enum:
            - glm-4.6
        messages:
          type: array
          description: >-
            The current conversation message list as the model’s prompt input,
            provided in JSON array format, e.g.,`{“role”: “user”, “content”:
            “Hello”}`. Possible message types include system messages, user
            messages. Note: The input must not consist of system or assistant
            messages only.
          items:
            oneOf:
              - title: User Message
                type: object
                properties:
                  role:
                    type: string
                    enum:
                      - user
                    description: Role of the message author
                    default: user
                  content:
                    oneOf:
                      - type: array
                        description: >-
                          Multimodal message content, supports text, images,
                          video, file
                        items:
                          $ref: '#/components/schemas/VisionMultimodalContentItem'
                      - type: string
                        description: >-
                          Text message content (can switch to multimodal message
                          above)
                        example: >-
                          What opportunities and challenges will the Chinese
                          large model industry face in 2025?
                required:
                  - role
                  - content
              - title: System Message
                type: object
                properties:
                  role:
                    type: string
                    enum:
                      - system
                    description: Role of the message author
                    default: system
                  content:
                    oneOf:
                      - type: string
                        description: Message text content
                        example: You are a helpful assistant.
                required:
                  - role
                  - content
              - title: Assistant Message
                type: object
                description: Can include tool calls
                properties:
                  role:
                    type: string
                    enum:
                      - assistant
                    description: Role of the message author
                    default: assistant
                  content:
                    oneOf:
                      - type: string
                        description: Text message content
                        example: I'll help you with that analysis.
                required:
                  - role
          minItems: 1
        tools:
          type: array
          description: List of tools the model can call. Supports up to `128` functions.
          anyOf:
            - items:
                $ref: '#/components/schemas/FunctionToolSchema'
        request_id:
          type: string
          description: >-
            Passed by the user side, needs to be unique; used to distinguish
            each request, 6–64 characters. If not provided by the user side, the
            platform will generate one by default.
          minLength: 6
          maxLength: 64
        user_id:
          type: string
          description: >-
            Unique ID for the end user, 6–128 characters. Avoid using sensitive
            information.
          minLength: 6
          maxLength: 128
    TokenizerResponse:
      type: object
      properties:
        created:
          type: integer
          format: int64
          example: 1727156815
        id:
          type: string
          example: 20241120141244890ab4ee4af84acf
          description: >-
            The task sequence number generated by the Zhipu AI Open Platform.
            Please use this number when calling the request result interface.
        request_id:
          type: string
          example: '1'
          description: >-
            The task number submitted by the client or generated by the platform
            when the request was initiated.
        usage:
          type: object
          properties:
            prompt_tokens:
              type: number
              description: Prompt tokens in this input
            image_tokens:
              type: number
              description: Image tokens in this input
            video_tokens:
              type: number
              description: Video tokens in this input
            total_tokens:
              type: number
              description: Total tokens in this input
      required:
        - id
        - usage
    Error:
      required:
        - code
        - message
      type: object
      description: The request has failed.
      properties:
        code:
          type: integer
          format: int32
          description: Error code.
        message:
          type: string
          description: Error message.
    VisionMultimodalContentItem:
      oneOf:
        - title: Text
          type: object
          properties:
            type:
              type: string
              enum:
                - text
              description: Content type is text
              default: text
            text:
              type: string
              description: Text content
          required:
            - type
            - text
          additionalProperties: false
        - title: Image
          type: object
          properties:
            type:
              type: string
              enum:
                - image_url
              description: Content type is image URL
              default: image_url
            image_url:
              type: object
              description: Image information
              properties:
                url:
                  type: string
                  description: >-
                    Image URL or Base64 encoding. Image size limit is under 5M
                    per image, with pixels not exceeding 6000*6000. GLM-5V are
                    limited to 150 sheets.Supports jpg, png, jpeg formats.
              required:
                - url
              additionalProperties: false
          required:
            - type
            - image_url
          additionalProperties: false
        - title: Video
          type: object
          properties:
            type:
              type: string
              enum:
                - video_url
              description: Content type is video URL
              default: video_url
            video_url:
              type: object
              description: Video information.
              properties:
                url:
                  type: string
                  description: >-
                    Video URL address.The video size is limited to within 200
                    MB,and the format supports `mp4`，`mkv`，`mov`.
              required:
                - url
              additionalProperties: false
          required:
            - type
            - video_url
          additionalProperties: false
        - title: File
          type: object
          properties:
            type:
              type: string
              enum:
                - file
              description: >-
                Content type is file. New unified file type, compatible with the
                legacy `file_url` type but not recommended for new scenarios.
                Not support passing both the `file` and `image_url` or
                `video_url` parameters at the same time.
              default: file
            file:
              type: object
              description: >-
                File content, supports one of `file_id`, `file_url`, or
                `file_data`. Single file size limit is `50M`.
              properties:
                file_id:
                  type: string
                  description: >-
                    The ID returned by the [File Upload
                    API](/api-reference/agents/file-upload), only GLM-5.3-Flash
                    series supported.
                file_url:
                  type: string
                  description: >-
                    File URL address. Only GLM-5.3-Flash series. Supports
                    formats such as pdf, txt, word, jsonl, xlsx, pptx, with a
                    maximum of 50.
                file_data:
                  type: string
                  description: >-
                    Base64 file content in the format
                    `data:<MIME>;base64,<BASE64_DATA>`.
                filename:
                  type: string
                  description: File name.
              additionalProperties: false
          required:
            - type
            - file
          additionalProperties: false
    FunctionToolSchema:
      type: object
      title: Function Call
      properties:
        type:
          type: string
          default: function
          enum:
            - function
        function:
          $ref: '#/components/schemas/FunctionObject'
      required:
        - type
        - function
      additionalProperties: false
    FunctionObject:
      type: object
      properties:
        name:
          type: string
          description: >-
            The name of the function to be called. Must be a-z, A-Z, 0-9, or
            contain underscores and dashes, with a maximum length of 64.
          minLength: 1
          maxLength: 64
          pattern: ^[a-zA-Z0-9_-]+$
        description:
          type: string
          description: >-
            A description of what the function does, used by the model to choose
            when and how to call the function.
        parameters:
          $ref: '#/components/schemas/FunctionParameters'
      required:
        - name
        - description
        - parameters
    FunctionParameters:
      type: object
      description: >-
        Parameters defined using JSON Schema. Must pass a JSON Schema object to
        accurately define accepted parameters. Omit if no parameters are needed
        when calling the function.
      additionalProperties: true
  securitySchemes:
    bearerAuth:
      type: http
      scheme: bearer
      description: >-
        Use the following format for authentication: Bearer [<your api
        key>](https://z.ai/manage-apikey/apikey-list)

````

This documentation is built and hosted on [Mintlify](https://mintlify.com), a developer documentation platform.