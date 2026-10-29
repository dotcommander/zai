> ## Documentation Index
> Fetch the complete documentation index at: https://docs.z.ai/llms.txt
> Use this file to discover all available pages before exploring further.

# Generate Video(Async)

> CogVideoX is a video generation large model developed by Z.AI, equipped with powerful video generation capabilities. Simply inputting text or images allows for effortless video creation.



## OpenAPI

````yaml POST /paas/v4/videos/generations
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
  /paas/v4/videos/generations:
    post:
      description: >-
        CogVideoX is a video generation large model developed by Z.AI, equipped
        with powerful video generation capabilities. Simply inputting text or
        images allows for effortless video creation.
      parameters:
        - $ref: '#/components/parameters/AcceptLanguage'
      requestBody:
        content:
          application/json:
            schema:
              oneOf:
                - $ref: '#/components/schemas/CogVideoX3Request'
                  title: CogVideoX-3
            examples:
              Text to Video Example:
                value:
                  model: cogvideox-3
                  prompt: A cat is playing with a ball.
                  quality: quality
                  with_audio: true
                  size: 1920x1080
                  fps: 30
              Image to Video Example:
                value:
                  model: cogvideox-3
                  image_url: >-
                    https://img.iplaysoft.com/wp-content/uploads/2019/free-images/free_stock_photo.jpg
                  prompt: Make the picture move
                  quality: quality
                  with_audio: true
                  size: 1920x1080
                  fps: 30
              First Last Frame to Video:
                value:
                  model: cogvideox-3
                  image_url:
                    - >-
                      https://gd-hbimg.huaban.com/ccee58d77afe8f5e17a572246b1994f7e027657fe9e6-qD66In_fw1200webp
                    - >-
                      https://gd-hbimg.huaban.com/cc2601d568a72d18d90b2cc7f1065b16b2d693f7fa3f7-hDAwNq_fw1200webp
                  prompt: Make the picture move
                  quality: quality
                  with_audio: true
                  size: 1920x1080
                  fps: 30
        required: true
      responses:
        '200':
          description: Processing successful.
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/VideoResponse'
        default:
          description: The request has failed.
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
components:
  parameters:
    AcceptLanguage:
      name: Accept-Language
      in: header
      schema:
        type: string
        description: Config desired response language for HTTP requests.
        default: en-US,en
        example: en-US,en
        enum:
          - en-US,en
      required: false
  schemas:
    CogVideoX3Request:
      allOf:
        - type: object
          properties:
            model:
              type: string
              description: The model code to be called.
              enum:
                - cogvideox-3
            prompt:
              type: string
              description: >-
                Text description of the video, maximum input length of 512
                characters. Either image_url or prompt must be provided, or
                both.
            quality:
              type: string
              description: >-
                Output mode, default is `speed`.

                - `quality`: Prioritizes quality, higher generation quality. 

                - `speed`: Prioritizes speed, faster generation time, relatively
                lower quality.
              example: speed
              enum:
                - speed
                - quality
            with_audio:
              type: boolean
              description: >-
                Whether to generate AI sound effects. Default: `false` (no sound
                effects).
              example: false
            image_url:
              type: array
              description: >-
                Provide an image based on which content will be generated. If
                this parameter is passed, the system will operate based on this
                image. Supports passing images via URL or Base64 encoding. Image
                requirements: images support `.png`, `.jpeg`, `.jpg` formats;
                image size: no more than `5M`. Either image_url and prompt can
                be used, or both can be passed simultaneously.

                First and last frames: supports inputting two images. The first
                uploaded image is regarded as the first frame, and the second
                image is regarded as the last frame. The model will generate the
                video based on the images passed in this parameter.

                First and last frame mode only supports `speed` mode
              items:
                oneOf:
                  - title: Image URL
                    type: string
                    format: uri
                    example: https://example.com/image.png
                  - title: Base64 Encoded Image
                    type: string
                    format: byte
                    example: data:image/png;base64, XXX
            size:
              type: string
              description: >-
                Default value: if not specified, the short side of the generated
                video is 1080 by default, and the long side is determined
                according to the original image ratio. Maximum support for 4K
                resolution. Resolution options: "1280x720", "720x1280",
                "1024x1024", "1080x1920", "2048x1080", "3840x2160"
              example: 1920x1080
              enum:
                - 1280x720
                - 720x1280
                - 1024x1024
                - 1920x1080
                - 1080x1920
                - 2048x1080
                - 3840x2160
            fps:
              type: integer
              description: >-
                Video frame rate (FPS), optional values are `30` or `60`.
                Default: `30`.
              example: 30
              enum:
                - 30
                - 60
            duration:
              type: integer
              description: >-
                Video duration, default is 5 seconds, supports `5` and `10`
                seconds.
              example: 5
              enum:
                - 5
                - 10
          required:
            - model
        - $ref: '#/components/schemas/VideoCommonRequest'
    VideoResponse:
      type: object
      properties:
        model:
          description: Model name used in this call.
          type: string
        id:
          description: >-
            Task order number generated by the Z.AI, use this order number when
            calling the request result interface.
          type: string
        request_id:
          description: >-
            Task number submitted by the user during the client request or
            generated by the platform.
          type: string
        task_status:
          description: >-
            Processing status, `PROCESSING (processing)`,` SUCCESS (success)`,
            `FAIL (failure)`. Results need to be obtained via query.
          type: string
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
    VideoCommonRequest:
      type: object
      properties:
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
            Unique ID of the end-user, assists the platform in intervening in
            end-user violations, generating illegal or inappropriate
            information, or other abusive behaviors. ID length requirement:
            minimum `6` characters, maximum `128` characters.
  securitySchemes:
    bearerAuth:
      type: http
      scheme: bearer
      description: >-
        Use the following format for authentication: Bearer [<your api
        key>](https://z.ai/manage-apikey/apikey-list)

````

This documentation is built and hosted on [Mintlify](https://mintlify.com), a developer documentation platform.