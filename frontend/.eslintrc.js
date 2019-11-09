module.exports = {
    env: {
        browser: true,
        es6: true,
    },
    extends: [
        "react-app",
        "eslint:recommended",
        "plugin:@typescript-eslint/eslint-recommended",
    ],
    globals: {
        Atomics: "readonly",
        SharedArrayBuffer: "readonly",
    },
    parser: "@typescript-eslint/parser",
    parserOptions: {
        ecmaFeatures: {
            jsx: true,
        },
        ecmaVersion: 2018,
        sourceType: "module",
    },
    plugins: ["react", "@typescript-eslint"],
    rules: {
        "react/jsx-indent": ["warn", 2],
        indent: ["warn", 4, { ignoredNodes: ["JSXElement"] }],
        quotes: "warn",
        semi: "warn",
        "comma-dangle": ["warn", "always-multiline"]
    },
};
