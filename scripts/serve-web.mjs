import { createServer } from "node:http";
import { readFileSync } from "node:fs";
import { fileURLToPath } from "node:url";
import path from "node:path";

const root = path.resolve(path.dirname(fileURLToPath(import.meta.url)), "..");
const source = readFileSync(path.join(root, "apps/web/src/main.tsx"), "utf8");

const readString = (key) => {
  const match = source.match(new RegExp(`${key}:\\s*"([^"]+)"`));
  if (!match) {
    throw new Error(`missing ${key} in apps/web/src/main.tsx`);
  }
  return match[1];
};

const shell = {
  testId: readString("testId"),
  title: readString("title"),
  status: readString("status"),
};

const html = `<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>${shell.title}</title>
  </head>
  <body>
    <main data-testid="${shell.testId}">
      <h1>${shell.title}</h1>
      <p>${shell.status}</p>
    </main>
  </body>
</html>`;

const server = createServer((req, res) => {
  if (req.url === "/" || req.url === "/index.html") {
    res.writeHead(200, { "content-type": "text/html; charset=utf-8" });
    res.end(html);
    return;
  }

  res.writeHead(404, { "content-type": "text/plain; charset=utf-8" });
  res.end("not found");
});

server.listen(4173, "127.0.0.1", () => {
  console.log("web server listening on http://127.0.0.1:4173");
});
