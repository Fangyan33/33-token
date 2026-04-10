import { spawnSync } from "node:child_process";
import path from "node:path";

const execDir = path.dirname(process.execPath);
const nodeVersionRoot = path.dirname(execDir);
const globalNodeModules = path.join(nodeVersionRoot, "lib", "node_modules");

const env = {
  ...process.env,
  NODE_PATH: [globalNodeModules, process.env.NODE_PATH].filter(Boolean).join(path.delimiter),
};

const result = spawnSync(
  "npm",
  [
    "exec",
    "--yes",
    "playwright",
    "test",
    "tests/e2e/platform-smoke.spec.ts",
    "--config=playwright.config.mjs",
  ],
  {
    stdio: "inherit",
    env,
  },
);

if (result.error) {
  throw result.error;
}

process.exit(result.status ?? 1);
