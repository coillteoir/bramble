module.exports = {
    root: true,
    env: { browser: true, es2020: true },
    extends: ["eslint:recommended", "plugin:@typescript-eslint/recommended"],
    ignorePatterns: ["dist", ".eslintrc.cjs", "src/bramble-types"],
    parser: "@typescript-eslint/parser",
    rules: {
        "arrow-body-style": ["warn"],
        "no-shadow": ["error"],
        "no-unneeded-ternary": ["warn"],
        "no-unreachable": ["error"],
        eqeqeq: ["error"],
        "max-depth": ["error", 4],
        "no-var": ["error"],
        "@typescript-eslint/no-explicit-any": ["warn"],
    },
};
