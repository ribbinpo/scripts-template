#!/usr/bin/env node

// Simple script to start the Pino Express server
const { spawn } = require('child_process');
const path = require('path');

console.log('🚀 Starting Pino Express Server...');
console.log('📊 Server will run on http://localhost:3002');
console.log('📝 Logs will be formatted using your existing Pino configuration');
console.log('');

// Start the server using tsx for TypeScript execution
const server = spawn('npx', ['tsx', 'src/server-pino.ts'], {
  cwd: __dirname,
  stdio: 'inherit',
  shell: true
});

server.on('error', (err) => {
  console.error('Failed to start server:', err);
});

server.on('close', (code) => {
  console.log(`Server process exited with code ${code}`);
});

// Handle graceful shutdown
process.on('SIGINT', () => {
  console.log('\n🛑 Shutting down server...');
  server.kill('SIGINT');
});

process.on('SIGTERM', () => {
  console.log('\n🛑 Shutting down server...');
  server.kill('SIGTERM');
});
