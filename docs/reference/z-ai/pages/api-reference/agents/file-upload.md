> ## Documentation Index
> Fetch the complete documentation index at: https://docs.z.ai/llms.txt
> Use this file to discover all available pages before exploring further.

# File Upload

> This API is designed for uploading auxiliary files (such as glossaries, terminology lists) to support the translation service. It allows users to upload reference materials that can enhance translation accuracy and consistency.

## File Limitations

* Maximum 100 files
* Maximum 100MB per file
* Files retained for 180 days


## OpenAPI

````yaml POST /paas/v4/files
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
  /paas/v4/files:
    post:
      description: >-
        This API is designed for uploading auxiliary files (such as glossaries,
        terminology lists) to support the translation service. It allows users
        to upload reference materials that can enhance translation accuracy and
        consistency.
      requestBody:
        content:
          multipart/form-data:
            schema:
              type: object
              properties:
                purpose:
                  type: string
                  description: Upload purpose
                  default: user_data
                  enum:
                    - user_data
                    - agent
                file:
                  type: string
                  format: binary
                  description: >-
                    File to upload. `agent` limit to `100MB`, `user_data` limit
                    to `1T`. `agent` allowed formats: `pdf`, `doc`, `xlsx`,
                    `ppt`, `txt`, `jpg`, `png`, `user_data` allowed formats:
                    `pptx`, `ppt`, `docx`, `doc`, `xlsx`, `xls`, `pdf`.
              required:
                - purpose
                - file
        required: true
      responses:
        '200':
          description: Processing successful
          content:
            application/json:
              schema:
                type: object
                properties:
                  id:
                    type: string
                    description: Unique identifier of the uploaded file.
                  object:
                    type: string
                    description: Object type.
                  bytes:
                    type: integer
                    description: File size in bytes.
                  filename:
                    type: string
                    description: Name of the uploaded file.
                  purpose:
                    type: string
                    description: Purpose of the uploaded file.
                  created_at:
                    type: integer
                    description: Timestamp of file creation.
        default:
          description: The request has failed.
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
components:
  schemas:
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