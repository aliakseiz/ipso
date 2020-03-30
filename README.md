Install openapi-code-generator:

```npm install @openapitools/openapi-generator-cli -g```

Generate client:

```npx openapi-generator generate -i api/openapi.yaml -g go -o api/client```

Client with XML tags:
```npx openapi-generator generate -i openapi.yaml -p withXml=true -g go -o client```