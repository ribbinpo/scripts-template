const tsConfigPaths = require('tsconfig-paths');
const path = require('path');

// Register the path mappings from tsconfig.json
tsConfigPaths.register({
  baseUrl: path.resolve(__dirname),
  paths: {
    "@/*": ["./dist/*"]
  }
}); 