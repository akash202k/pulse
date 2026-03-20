// dummy-multi-server.js
const http = require("http");

const BASE_PORT = 3000;
const SERVER_COUNT = 20; // change to 10 if needed

function rand(min, max) {
  return Math.floor(Math.random() * (max - min + 1)) + min;
}

function startServer(port) {
  const server = http.createServer((req, res) => {
    const match = req.url.match(/^\/api\/app(\d+)\/health$/);

    if (!match) {
      res.writeHead(404);
      return res.end("Not Found");
    }

    const appId = match[1];

    const latency = rand(100, 2000);
    const fail = Math.random() < 0.2;
    const hang = Math.random() < 0.05;

    if (hang) {
      console.log(`[${port}] app${appId} 🟡 HANG`);
      return;
    }

    setTimeout(() => {
      if (fail) {
        console.log(`[${port}] app${appId} ❌`);
        res.writeHead(500);
        return res.end("DOWN");
      }

      console.log(`[${port}] app${appId} ✅`);
      res.writeHead(200);
      res.end("UP");
    }, latency);
  });

  server.listen(port, () => {
    console.log(`Dummy server running on http://localhost:${port}`);
  });
}

// Start N servers
for (let i = 0; i < SERVER_COUNT; i++) {
  startServer(BASE_PORT + i);
}
