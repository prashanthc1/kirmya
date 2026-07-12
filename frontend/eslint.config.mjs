import { defineConfig, globalIgnores } from "eslint/config";
import nextVitals from "eslint-config-next/core-web-vitals";
import nextTs from "eslint-config-next/typescript";

const eslintConfig = defineConfig([
  ...nextVitals,
  ...nextTs,
  {
    rules: {
      // React-Compiler guidance rules (eslint-plugin-react-hooks v6). These flag
      // idiomatic patterns this app uses deliberately — fetch-on-mount effects
      // and Date.now() inside event handlers — so keep them as warnings (real
      // issues still surface) rather than failing the build/CI on style.
      "react-hooks/set-state-in-effect": "warn",
      "react-hooks/purity": "warn",
      // Same rationale: this rule fires on the app's use-before-declare
      // fetch-on-mount effect pattern. Keep it visible but non-blocking.
      "react-hooks/immutability": "warn",
      // The app leans on `any` at several API/adapter boundaries; surface it as
      // a warning rather than a hard CI failure (matches the style rules above).
      "@typescript-eslint/no-explicit-any": "warn",
    },
  },
  // Override default ignores of eslint-config-next.
  globalIgnores([
    // Default ignores of eslint-config-next:
    ".next/**",
    "out/**",
    "build/**",
    "next-env.d.ts",
  ]),
]);

export default eslintConfig;
