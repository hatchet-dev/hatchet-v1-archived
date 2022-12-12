#!/bin/bash

swagger generate spec --scan-models --include github.com/hatchet-dev/hatchet --work-dir ./cmd/hatchet-server/ --output ./bin/oas/api-server-generated.yaml

# add servers block to generated
echo 'servers:
  - url: "http://localhost:8080"
components:
  securitySchemes:
    Bearer:
      type: http
      scheme: bearer
security:
- Bearer: []
' > ./bin/oas/api-server-generated-extra.yaml

yq eval-all 'select(fi == 0) * select(filename == "./bin/oas/api-server-generated.yaml")' ./bin/oas/api-server-generated-extra.yaml ./bin/oas/api-server-generated.yaml > ./bin/oas/api-server-generated-with-mixins.yaml

node ./hack/oas/samples-generator/index.js ./bin/oas/api-server-generated-with-mixins.yaml ./bin/oas/api-server-generated-with-samples.yaml

# convert the server-generated API file to 3.0
openapi-generator-cli generate -g openapi-yaml -i ./bin/oas/api-server-generated-with-samples.yaml -o ./api-server-3.0 --skip-validate-spec

# merge the base api with the server 3.0 file
yq eval-all 'select(fi == 0) * select(filename == "./hack/oas/api-base.yaml")' api-server-3.0/openapi/openapi.yaml ./hack/oas/api-base.yaml > ./bin/oas/api-populated-pre-fmt.yaml

cat ./bin/oas/api-populated-pre-fmt.yaml | sed 's/%7B/{/g' | sed 's/%7D/}/g' | sed 's/\*\/\*/application\/json/g' > ./bin/oas/openapi.yaml

# remove the generated directories
rm -rf ./api-server-3.0
rm openapitools.json
rm ./bin/oas/api-populated-pre-fmt.yaml
rm ./bin/oas/api-server-generated-with-samples.yaml
rm ./bin/oas/api-server-generated-with-mixins.yaml
rm ./bin/oas/api-server-generated-extra.yaml
rm ./bin/oas/api-server-generated.yaml