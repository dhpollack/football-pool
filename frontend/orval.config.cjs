module.exports = {
  "football-pool-api": {
    input: "../backend/openapi.json",
    output: {
      mock: true,
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
