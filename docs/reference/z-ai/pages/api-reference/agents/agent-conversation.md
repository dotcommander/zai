> ## Documentation Index
> Fetch the complete documentation index at: https://docs.z.ai/llms.txt
> Use this file to discover all available pages before exploring further.

# Conversation History

> This endpoint is used to query the agent conversation history.Only support slides_glm_agent



## OpenAPI

````yaml POST /v1/agents/conversation
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
  /v1/agents/conversation:
    post:
      description: >-
        This endpoint is used to query the agent conversation history.Only
        support slides_glm_agent
      parameters:
        - $ref: '#/components/parameters/AcceptLanguage'
      requestBody:
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/GlmSlideAgentConversationRequest'
        required: true
      responses:
        '200':
          description: Processing successful
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/GlmSlideAgentConversationResponse'
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
    GlmSlideAgentConversationRequest:
      type: object
      properties:
        agent_id:
          type: string
          description: Agent ID
        conversation_id:
          type: string
          description: Conversation ID
        custom_variables:
          type: object
          description: Custom variables
          properties:
            include_pdf:
              type: boolean
              description: Is export the pdf file
            pages:
              type: array
              description: Slides Pages
              items:
                type: object
                properties:
                  position:
                    type: number
                    description: Slide Page Position
                  width:
                    type: number
                    description: 'Slide Width, unit: pt'
                  height:
                    type: number
                    description: 'Slide Height, unit: pt'
    GlmSlideAgentConversationResponse:
      type: object
      properties:
        conversation_id:
          type: string
          description: Conversation ID
        agent_id:
          type: string
          description: Agent ID
        choices:
          type: array
          description: Agent output.
          items:
            type: object
            properties:
              message:
                type: array
                items:
                  type: object
                  properties:
                    role:
                      type: string
                      description: 'Role: fixed as `assistant`.'
                    content:
                      type: array
                      description: Content metadata
                      items:
                        type: object
                        properties:
                          type:
                            type: string
                            description: 'Response Content type: file_url、image_url'
                          tag_cn:
                            type: string
                            description: CN Tag.
                          tag_en:
                            type: string
                            description: EN Tag.
                          file_url:
                            type: string
                            description: Output file_url content when type is file_url
                          image_url:
                            type: string
                            description: Output image_url content when type is image_url
        error:
          type: object
          properties:
            code:
              type: string
              description: Error code.
            message:
              type: string
              description: Error message.
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