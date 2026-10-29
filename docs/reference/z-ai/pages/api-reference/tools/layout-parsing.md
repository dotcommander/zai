> ## Documentation Index
> Fetch the complete documentation index at: https://docs.z.ai/llms.txt
> Use this file to discover all available pages before exploring further.

# Layout Parsing

> Use the [GLM-OCR](/guides/vlm/glm-ocr) model to parse the layout of documents and images and extract text content. Support OCR recognition of images and PDF documents, returning detailed layout information and visualization results.



## OpenAPI

````yaml POST /paas/v4/layout_parsing
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
  /paas/v4/layout_parsing:
    post:
      summary: Layout Parsing
      description: >-
        Use the [GLM-OCR](/guides/vlm/glm-ocr) model to parse the layout of
        documents and images and extract text content. Support OCR recognition
        of images and PDF documents, returning detailed layout information and
        visualization results.
      requestBody:
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/LayoutParsingRequest'
            examples:
              Layout Parsing Example:
                value:
                  model: glm-ocr
                  file: https://cdn.bigmodel.cn/static/logo/introduction.png
        required: true
      responses:
        '200':
          description: Business processing successful
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/LayoutParsingResponse'
        default:
          description: Request failed.
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
components:
  schemas:
    LayoutParsingRequest:
      type: object
      required:
        - model
        - file
      properties:
        model:
          type: string
          description: 'Model code: `glm-ocr`'
          example: glm-ocr
          enum:
            - glm-ocr
        file:
          type: string
          description: >-
            Image or PDF document to be recognized, supports URL and base64.
            Supported image formats: PDF, JPG, PNG. Single image ≤10MB, PDF
            ≤50MB, maximum support 30 pages
          example: https://cdn.bigmodel.cn/static/logo/introduction.png
        return_crop_images:
          type: boolean
          description: Whether to return screenshot information
          default: false
        need_layout_visualization:
          type: boolean
          description: Whether to return detailed layout image result information
          default: false
        start_page_id:
          type: integer
          description: Start page number for parsing when PDF is provided
          minimum: 1
        end_page_id:
          type: integer
          description: End page number for parsing when PDF is provided
          minimum: 1
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
    LayoutParsingResponse:
      type: object
      properties:
        id:
          type: string
          description: Task ID
          example: task_123456789
        created:
          type: integer
          format: int64
          description: Request creation time, Unix timestamp in seconds
          example: 1727156815
        model:
          type: string
          description: Model name
          example: GLM-OCR
        md_results:
          type: string
          description: Recognition result in Markdown format
          example: |-
            # Doc title
            This is the document content...
        layout_details:
          type: array
          description: Detailed layout information
          items:
            type: array
            items:
              $ref: '#/components/schemas/LayoutDetail'
        layout_visualization:
          type: array
          description: Recognition result image URLs
          items:
            type: string
        data_info:
          $ref: '#/components/schemas/DataInfo'
        usage:
          type: object
          description: Token usage statistics returned when the model call ends.
          properties:
            prompt_tokens:
              type: number
              description: Number of tokens in user input
            completion_tokens:
              type: number
              description: Number of output tokens
            prompt_tokens_details:
              type: object
              properties:
                cached_tokens:
                  type: number
                  description: Number of tokens served from cache
            total_tokens:
              type: integer
              description: Total number of tokens
        request_id:
          type: string
          description: Request ID
          example: req_123456789
      required:
        - id
        - created
        - model
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
    LayoutDetail:
      type: object
      description: Layout detail element
      properties:
        index:
          type: integer
          description: Element index
          example: 1
        label:
          type: string
          description: >-
            Element type: image for images, text for text content, formula for
            inline formulas, table for tables
          enum:
            - image
            - text
            - formula
            - table
          example: text
        bbox_2d:
          type: array
          description: Normalized element coordinates [x1,y1,x2,y2]
          items:
            type: number
            minimum: 0
            maximum: 1
          minItems: 4
          maxItems: 4
          example:
            - 0.1
            - 0.1
            - 0.5
            - 0.3
        content:
          type: string
          description: Element content (text / image URL / table HTML)
          example: This is the content of the element
        height:
          type: integer
          description: Page height
          example: 800
        width:
          type: integer
          description: Page width
          example: 600
      required:
        - index
        - label
    DataInfo:
      type: object
      description: Document basic information
      properties:
        num_pages:
          type: integer
          description: Total number of document pages
          example: 5
        pages:
          type: array
          description: Document page count information
          items:
            $ref: '#/components/schemas/PageInfo'
      required:
        - num_pages
    PageInfo:
      type: object
      description: Page dimension information
      properties:
        width:
          type: integer
          description: Page width
          example: 600
        height:
          type: integer
          description: Page height
          example: 800
      required:
        - width
        - height
  securitySchemes:
    bearerAuth:
      type: http
      scheme: bearer
      description: >-
        Use the following format for authentication: Bearer [<your api
        key>](https://z.ai/manage-apikey/apikey-list)

````

This documentation is built and hosted on [Mintlify](https://mintlify.com), a developer documentation platform.