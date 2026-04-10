const { execSync } = require("node:child_process");
const { createRequire } = require("node:module");

const globalRoot = execSync("npm root -g", { encoding: "utf8" }).trim();
const requireFromGlobal = createRequire(`${globalRoot}/playwright/package.json`);

module.exports = requireFromGlobal("playwright/test");
