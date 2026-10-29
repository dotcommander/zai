> ## Documentation Index
> Fetch the complete documentation index at: https://docs.z.ai/llms.txt
> Use this file to discover all available pages before exploring further.

# Web Search

> The [Web Search](/guides/tools/web-search) is a specialized search engine for large language models. Building upon traditional search engine capabilities like web crawling and ranking, it enhances intent recognition to return results better suited for LLM processing (including webpage titles, URLs, summaries, site names, favicons etc.).



## OpenAPI

````yaml POST /paas/v4/web_search
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
  /paas/v4/web_search:
    post:
      description: >-
        The [Web Search](/guides/tools/web-search) is a specialized search
        engine for large language models. Building upon traditional search
        engine capabilities like web crawling and ranking, it enhances intent
        recognition to return results better suited for LLM processing
        (including webpage titles, URLs, summaries, site names, favicons etc.).
      parameters:
        - $ref: '#/components/parameters/AcceptLanguage'
      requestBody:
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/WebSearchRequest'
        required: true
      responses:
        '200':
          description: Processing successful
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/WebSearchResponse'
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
    WebSearchRequest:
      type: object
      properties:
        search_engine:
          type: string
          description: |-
            The search engine code to call.
             search-prime: Z.AI Premium Version Search Engine
          example: search-prime
          default: search-prime
          enum:
            - search-prime
        search_query:
          type: string
          description: The content to be searched.
        count:
          type: integer
          description: |-
            The number of results to return
            Fillable range: `1-50`, maximum `50` results per single search
            Default is `10`
            Supported search engines: 
            `search_pro_jina`.
          minimum: 1
          maximum: 50
        search_domain_filter:
          type: string
          description: >-
            Used to limit the scope of search results and only return content
            from specified whitelist domains.

            Whitelist: Directly enter the domain name (e.g., `www.example.com`)

            Supported search engines: 

            `search_pro_jina`
        search_recency_filter:
          type: string
          description: |-
            Search for webpages within a specified time range.
            Default is `noLimit`
            Fillable values:
            `oneDay`: within one day
            `oneWeek`: within one week
            `oneMonth`: within one month
            `oneYear`: within one year
            `noLimit`: no limit (default)
            Supported search engines: 
            `search_pro_jina`
          enum:
            - oneDay
            - oneWeek
            - oneMonth
            - oneYear
            - noLimit
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
      required:
        - search_engine
        - search_query
    WebSearchResponse:
      type: object
      properties:
        id:
          description: Task ID.
          type: string
        created:
          description: Request creation time, Unix timestamp in seconds.
          type: integer
        search_result:
          description: Search results.
          type: array
          items:
            $ref: '#/components/schemas/WebSearchObjectResponse'
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
    WebSearchObjectResponse:
      type: object
      properties:
        title:
          type: string
          description: Title.
        content:
          type: string
          description: Content summary.
        link:
          type: string
          description: Result URL.
        media:
          type: string
          description: Website name.
        icon:
          type: string
          description: Website icon.
        refer:
          type: string
          description: Index number.
        publish_date:
          type: string
          description: Website publication date.
  securitySchemes:
    bearerAuth:
      type: http
      scheme: bearer
      description: >-
        Use the following format for authentication: Bearer [<your api
        key>](https://z.ai/manage-apikey/apikey-list)

````

This documentation is built and hosted on [Mintlify](https://mintlify.com), a developer documentation platform.