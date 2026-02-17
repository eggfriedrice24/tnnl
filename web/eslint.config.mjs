import antfu from "@antfu/eslint-config";

export default antfu(
  {
    type: "app",
    react: true,
    typescript: true,
    formatters: true,
    stylistic: {
      indent: 2,
      semi: true,
      quotes: "double",
    },
    ignores: [
      "node_modules/*",
      "dist/*",
      ".claude/**",
      "**/*.md",
      "src/routeTree.gen.ts",
    ],
  },
  {
    rules: {
      "react-refresh/only-export-components": "off",
      "ts/consistent-type-definitions": ["error", "type"],
      "no-console": "warn",
      "antfu/no-top-level-await": "off",
      "node/prefer-global/process": "off",
      "perfectionist/sort-imports": ["error", {
        groups: [
          "builtin",
          "react",
          "external",
          "internal",
          ["parent", "sibling", "index"],
          "side-effect",
          "unknown",
        ],
        customGroups: [
          { groupName: "react", elementNamePattern: "^react(-.*)?$|^react-dom$" },
        ],
        internalPattern: ["^@tnnl/.+", "^@/.+"],
        newlinesBetween: 1,
      }],
      "unicorn/filename-case": [
        "error",
        {
          case: "kebabCase",
          ignore: ["README.md", "CLAUDE.md"],
        },
      ],
    },
  },
);
