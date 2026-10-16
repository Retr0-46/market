const http = require('node:http');
const fs = require('node:fs');
const path = require('node:path');

const port = Number(process.env.FRONTEND_PORT || 3000);
const productService = process.env.PRODUCT_SERVICE_URL || 'http://localhost:8081';
const authService = process.env.AUTH_SERVICE_URL || 'http://localhost:8080';
const cartService = process.env.CART_SERVICE_URL || 'http://localhost:8082';
const files = {
  '/': ['index.html', 'text/html; charset=utf-8'],
  '/app.js': ['app.js', 'text/javascript; charset=utf-8'],
  '/style.css': ['style.css', 'text/css; charset=utf-8']
};

http.createServer((request, response) => {
  const pathname = new URL(request.url, 'http://localhost').pathname;

  const service = pathname === '/products' || pathname.startsWith('/products/')
    ? productService
    : pathname === '/auth' || pathname.startsWith('/auth/')
      ? authService
      : pathname === '/cart' || pathname.startsWith('/cart/')
        ? cartService
        : null;

  if (service) {
    const target = new URL(request.url, service);
    const upstream = http.request(target, {
      method: request.method,
      headers: { ...request.headers, host: target.host }
    }, upstreamResponse => {
      response.writeHead(upstreamResponse.statusCode, upstreamResponse.headers);
      upstreamResponse.pipe(response);
    });
    upstream.on('error', () => {
      if (!response.headersSent) {
        response.writeHead(502, { 'Content-Type': 'text/plain; charset=utf-8' });
      }
      response.end('Service is unavailable');
    });
    request.pipe(upstream);
    return;
  }

  const file = files[pathname];
  if (!file || request.method !== 'GET') {
    response.writeHead(404);
    response.end();
    return;
  }

  response.writeHead(200, { 'Content-Type': file[1] });
  fs.createReadStream(path.join(__dirname, file[0])).pipe(response);
}).listen(port, () => {
  console.log(`Frontend: http://localhost:${port}`);
});
