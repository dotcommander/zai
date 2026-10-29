> ## Documentation Index
> Fetch the complete documentation index at: https://docs.z.ai/llms.txt
> Use this file to discover all available pages before exploring further.

# Generate Image

> Use [GLM-Image](/guides/image/glm-image) series models to generate high-quality images from text prompts. Through quick and accurate understanding of user text descriptions, `AI` image expression becomes more precise and personalized.



## OpenAPI

````yaml POST /paas/v4/images/generations
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
  /paas/v4/images/generations:
    post:
      summary: Generate Image
      description: >-
        Use [GLM-Image](/guides/image/glm-image) series models to generate
        high-quality images from text prompts. Through quick and accurate
        understanding of user text descriptions, `AI` image expression becomes
        more precise and personalized.
      requestBody:
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/CreateImageRequest'
            examples:
              Generate Image Example:
                value:
                  model: glm-image
                  prompt: >-
                    A cute little kitten sitting on a sunny windowsill, with the
                    background of blue sky and white clouds.
                  size: 1280x1280
        required: true
      responses:
        '200':
          description: Processing successful
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ImageGenerationResponse'
        default:
          description: Request Failed
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
components:
  schemas:
    CreateImageRequest:
      type: object
      required:
        - model
        - prompt
      properties:
        model:
          type: string
          description: Model code
          enum:
            - glm-image
          example: glm-image
        prompt:
          type: string
          description: The text description of the image to be generated.
          example: A cute little kitten.
        quality:
          type: string
          description: >-
            The quality of the generated image. `glm-image` default is `hd`,
            others model is `standard`. `hd`: Generates a more detailed and rich
            image with higher overall consistency, but takes about `20` seconds.
            `standard`: Generates an image quickly, suitable for scenarios with
            higher requirements for generation speed, takes about `5-10`
            seconds.
          enum:
            - hd
            - standard
          default: hd
        size:
          type: string
          description: >-
            Image size. `glm-image` recommended enum values: `1280x1280`
            (default), `1568x1056`, `1056x1568`, `1472x1088`, `1088x1472`,
            `1728x960`, `960x1728`. Custom parameter: Both width and height must
            be between `1024px-2048px`, and must be divisible by `32`, and the
            maximum pixel count must not exceed `2^22px`. 

            Others model recommended enum values: `1024x1024` (default),
            `768x1344`, `864x1152`, `1344x768`, `1152x864`, `1440x720`,
            `720x1440`. Custom parameter: Both width and height must be between
            `512px-2048px`, and must be divisible by `16`, and the maximum pixel
            count must not exceed `2^21px`.
          default: 1280x1280
          example: 1280x1280
        user_id:
          type: string
          description: >-
            Unique ID of the end user, helping the platform intervene in illegal
            activities, inappropriate content generation, or other abuses. ID
            length: 6 to 128 characters.
          minLength: 6
          maxLength: 128
    ImageGenerationResponse:
      type: object
      properties:
        created:
          type: integer
          example: 1760335349
          description: Request creation time, in `Unix` timestamp format, unit is seconds.
        data:
          type: array
          description: >-
            Array, containing the generated image `URL`. Currently, the array
            only contains one image.
          items:
            type: object
            properties:
              url:
                type: string
                description: >-
                  Image link. The temporary link expires after `30` days, please
                  store it promptly.
            required:
              - url
        content_filter:
          type: array
          description: Array, containing content safety related information.
          items:
            type: object
            properties:
              role:
                type: string
                description: >-
                  Safety enforcement stage, including `role = assistant` model
                  inference, `role = user` user input, `role = history`
                  historical context.
                enum:
                  - assistant
                  - user
                  - history
              level:
                type: integer
                description: >-
                  Severity level `level 0-3`, `level 0` is most severe, `3` is
                  least severe.
                minimum: 0
                maximum: 3
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