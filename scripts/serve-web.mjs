import { createServer } from "node:http";
import { createReadStream, existsSync, statSync } from "node:fs";
import { fileURLToPath } from "node:url";
import path from "node:path";
import { extname, normalize } from "node:path";

const root = path.resolve(path.dirname(fileURLToPath(import.meta.url)), "..");
const webRoot = path.join(root, "apps/web");

const mimeTypes = {
  ".html": "text/html; charset=utf-8",
  ".js": "text/javascript; charset=utf-8",
  ".mjs": "text/javascript; charset=utf-8",
  ".ts": "text/javascript; charset=utf-8",
  ".tsx": "text/javascript; charset=utf-8",
  ".css": "text/css; charset=utf-8",
};

const server = createServer((req, res) => {
  const url = req.url || "/";
  const cleanPath = normalize(url.split("?")[0]).replace(/^(\.\.(\/|\\|$))+/, "");
  const relativePath = cleanPath === "/" ? "index.html" : cleanPath.replace(/^\//, "");
  const filePath = path.join(webRoot, relativePath);

  if (!filePath.startsWith(webRoot) || !existsSync(filePath) || !statSync(filePath).isFile()) {
    res.writeHead(404, { "content-type": "text/plain; charset=utf-8" });
    res.end("not found");
    return;
  }

  const contentType = mimeTypes[extname(filePath)] || "application/octet-stream";
  res.writeHead(200, { "content-type": contentType });
  createReadStream(filePath).pipe(res);
});

server.listen(4173, "127.0.0.1", () => {
  console.log("web server listening on http://127.0.0.1:4173");
});
