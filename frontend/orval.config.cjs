module.exports = {
  "football-pool-api": {
    input: "../openapi.json",
    output: {
      mock: {
        type: "msw",
        delay: 10,
        generateEachHttpStatus: true,
        indexMockFiles: true,
      },
      mode: "tags-split",
      target: "./src/services/api",
      schemas: "./src/services/model",
      client: "react-query",
      override: {
        mutator: {
          path: "./src/services/custom-instance.ts",
          name: "customInstance",
        },
      },
    },
    hooks: {
      afterAllFilesWrite: "biome format --write ./src/services",
    },
  },
};
