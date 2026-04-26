import { createReadStream, existsSync } from "node:fs";
import { stat } from "node:fs/promises";
import http from "node:http";
import path from "node:path";
import { fileURLToPath } from "node:url";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
const host = "127.0.0.1";
const port = Number(process.env.PORT || 3000);
const backendUrl = "http://localhost:8080";

const contentTypes = {
  ".css": "text/css; charset=utf-8",
  ".html": "text/html; charset=utf-8",
  ".js": "application/javascript; charset=utf-8",
  ".json": "application/json; charset=utf-8",
};

const sendFile = async (response, filePath) => {
  try {
    const fileStat = await stat(filePath);
    if (!fileStat.isFile()) {
      response.writeHead(404);
      response.end("Not found");
      return;
    }
    const ext = path.extname(filePath);
    response.writeHead(200, {
      "Content-Type": contentTypes[ext] || "application/octet-stream",
      "Cache-Control": "no-cache",
    });
    createReadStream(filePath).pipe(response);
  } catch {
    response.writeHead(404);
    response.end("Not found");
  }
};

const proxyRequest = async (request, response, url) => {
  const headers = new Headers();

  Object.entries(request.headers).forEach(([key, value]) => {
    if (value !== undefined && key.toLowerCase() !== "host") {
      headers.set(key, value);
    }
  });

  const hasBody = request.method && !["GET", "HEAD"].includes(request.method.toUpperCase());
  const bodyBuffer = hasBody
    ? await new Promise((resolve, reject) => {
        const chunks = [];
        request.on("data", (chunk) => chunks.push(chunk));
        request.on("end", () => resolve(Buffer.concat(chunks)));
        request.on("error", reject);
      })
    : undefined;

  const upstreamResponse = await fetch(`${backendUrl}${url.pathname}${url.search}`, {
    method: request.method,
    headers,
    body: bodyBuffer,
    duplex: hasBody ? "half" : undefined,
  });

  response.writeHead(upstreamResponse.status, Object.fromEntries(upstreamResponse.headers.entries()));
  const arrayBuffer = await upstreamResponse.arrayBuffer();
  response.end(Buffer.from(arrayBuffer));
};

http
  .createServer(async (request, response) => {
    const url = new URL(request.url || "/", `http://${request.headers.host}`);

    if (url.pathname.startsWith("/api/")) {
      try {
        await proxyRequest(request, response, url);
      } catch (error) {
        response.writeHead(502, { "Content-Type": "application/json; charset=utf-8" });
        response.end(
          JSON.stringify({
            error: {
              code: "BAD_GATEWAY",
              message: "Не удалось связаться с backend-сервером",
              details: String(error),
            },
          }),
        );
      }
      return;
    }

    const requestedPath = path.join(__dirname, decodeURIComponent(url.pathname));

    if (existsSync(requestedPath) && path.extname(requestedPath)) {
      await sendFile(response, requestedPath);
      return;
    }

    await sendFile(response, path.join(__dirname, "index.html"));
  })
  .listen(port, host, () => {
    console.log(`Frontend is running at http://${host}:${port}`);
  })
  .on("error", (error) => {
    console.error(`Failed to start frontend server on http://${host}:${port}`);
    console.error(error);
  });
