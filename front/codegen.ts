
import type { CodegenConfig } from '@graphql-codegen/cli';

const config: CodegenConfig = {
  overwrite: true,
  schema: "../server/handler/graph/resolver/*.graphql",
  documents: "src/**/*.tsx",
  ignoreNoDocuments: true,
  generates: {
    "src/gql/": {
      preset: "client",
      plugins: [],
      config: {
        useTypeImports: true,
      }
    }
  },
};

export default config;
