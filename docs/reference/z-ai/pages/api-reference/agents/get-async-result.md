> ## Documentation Index
> Fetch the complete documentation index at: https://docs.z.ai/llms.txt
> Use this file to discover all available pages before exploring further.

# Retrieve Result

> This endpoint is used to query the result of an asynchronous request.



## OpenAPI

````yaml POST /v1/agents/async-result
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
  /v1/agents/async-result:
    post:
      description: This endpoint is used to query the result of an asynchronous request.
      parameters:
        - $ref: '#/components/parameters/AcceptLanguage'
      requestBody:
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/CommonAgentResultRequest'
        required: true
      responses:
        '200':
          description: Processing successful
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/CommonAgentResultResponse'
        default:
          description: The request has failed.
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/SpecialEffectsVideosAgentError'
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
    CommonAgentResultRequest:
      type: object
      properties:
        agent_id:
          type: string
          description: 'Agent ID: `vidu_template_agent`.'
        async_id:
          type: string
          description: Task ID from async response.
    CommonAgentResultResponse:
      type: object
      properties:
        status:
          type: string
          description: '`pending` (processing), `success` (completed), `failed` (failed).'
        agent_id:
          type: string
          description: Agent ID
        async_id:
          type: string
          description: Asynchronous task ID.
        choices:
          type: array
          description: Agent output.
          items:
            type: object
            properties:
              index:
                type: integer
                description: Result index.
              finish_reason:
                type: string
                description: >-
                  Reason for model inference termination. Can be ‘stop’,
                  ‘tool_calls’, ‘length’, ‘sensitive’, or ‘network_error’.
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
                      description: Video file metadata
                      items:
                        type: object
                        properties:
                          type:
                            type: string
                            description: 'object type: `video_url`.'
                          video_url:
                            type: string
                            description: MP4 video URL.
    SpecialEffectsVideosAgentError:
      type: object
      properties:
        status:
          type: string
          description: 'Status: `pending` (task created), `failed` (task creation failed).'
        agent_id:
          type: string
          description: Agent ID
        error:
          type: object
          properties:
            code:
              type: string
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