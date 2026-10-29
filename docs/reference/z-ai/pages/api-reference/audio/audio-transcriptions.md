> ## Documentation Index
> Fetch the complete documentation index at: https://docs.z.ai/llms.txt
> Use this file to discover all available pages before exploring further.

# Audio Transcriptions

> Use the [GLM-ASR-2512](/guides/audio/glm-asr-2512) model to transcribe audio files into text, supporting multiple languages and real-time streaming transcription.



## OpenAPI

````yaml POST /paas/v4/audio/transcriptions
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
  /paas/v4/audio/transcriptions:
    post:
      summary: Speech to Text
      description: >-
        Use the [GLM-ASR-2512](/guides/audio/glm-asr-2512) model to transcribe
        audio files into text, supporting multiple languages and real-time
        streaming transcription.
      requestBody:
        content:
          multipart/form-data:
            schema:
              $ref: '#/components/schemas/AudioTranscriptionRequest'
            example:
              model: glm-asr-2512
              stream: false
        required: true
      responses:
        '200':
          description: Request processed successfully
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/AudioTranscriptionResponse'
            text/event-stream:
              schema:
                $ref: '#/components/schemas/AudioTranscriptionStreamResponse'
        default:
          description: Request failed.
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
components:
  schemas:
    AudioTranscriptionRequest:
      type: object
      required:
        - file
        - model
      properties:
        file:
          type: string
          format: binary
          description: >-
            The audio file to be transcribed. Supported audio file formats:
            `.wav / .mp3`. Specifications: file size ≤ `25 MB`, audio duration ≤
            `30 seconds`.
        file_base64:
          type: string
          description: >-
            Base64 encoded audio file. Only one of file_base64 or file needs to
            be provided (if both are provided, file takes precedence).
        model:
          type: string
          description: The model ID to invoke.
          enum:
            - glm-asr-2512
          default: glm-asr-2512
        prompt:
          type: string
          description: >-
            In long text scenarios, you can provide previous transcription
            results as context. Recommended to be less than 8000 characters.
        hotwords:
          type: array
          description: >-
            Hotword list to improve recognition accuracy for domain-specific
            vocabulary. Format example: ["person_name","place_name"].
            Recommended not to exceed 100 items.
          items:
            type: string
          maxItems: 100
        stream:
          type: boolean
          default: false
          description: >-
            This parameter should be set to `false` or omitted when using
            synchronous calls. It indicates that the model returns all content
            at once after generating all content. Default is `false`. If set to
            `true`, the model will return generated content in chunks via
            standard `Event Stream`. When the `Event Stream` ends, a `data:
            [DONE]` message will be returned.
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
    AudioTranscriptionResponse:
      type: object
      properties:
        id:
          type: string
          description: Task ID
        created:
          type: integer
          format: int64
          description: Request creation time, as a `Unix` timestamp in seconds.
        request_id:
          type: string
          description: >-
            Passed by the client, must be unique. A unique identifier to
            distinguish each request. If not provided by the client, the
            platform will generate one by default.
        model:
          description: Model name
          type: string
        text:
          type: string
          description: The complete transcribed content of the audio.
    AudioTranscriptionStreamResponse:
      type: object
      properties:
        id:
          type: string
          description: Task ID
        created:
          type: integer
          format: int64
          description: Request creation time, as a `Unix` timestamp in seconds.
        model:
          description: Model name
          type: string
        type:
          type: string
          description: >-
            Audio transcription event type. `transcript.text.delta` indicates
            transcription in progress, `transcript.text.done` indicates
            transcription completed.
        delta:
          type: string
          description: Incremental audio transcription information returned by the model.
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
  securitySchemes:
    bearerAuth:
      type: http
      scheme: bearer
      description: >-
        Use the following format for authentication: Bearer [<your api
        key>](https://z.ai/manage-apikey/apikey-list)

````

This documentation is built and hosted on [Mintlify](https://mintlify.com), a developer documentation platform.