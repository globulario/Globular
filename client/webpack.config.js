/*
module.exports = {
  mode: "development",
  target: 'web',
  entry: "./lib/services.js",
  output: {
    filename: './lib/bundle.js'
  }
};
*/
/*
// Typescript code generation.
const isProduction = false;
const path = require('path');

module.exports = {
  entry: './services.ts',
  target: 'web',
  module: {
    rules: [
      {
        test: /\.tsx?$/,
        use: 'ts-loader',
        exclude: /node_modules/,
        options: {
          compilerOptions: {
            sourceMap: !isProduction,
          },
        },
      },
    ],
  },
  resolve: {
    extensions: ['.tsx', '.ts', '.js'],
  },
  output: {
    path: path.resolve(__dirname, 'dist'),
    filename: 'bundle.js',
  },
};
*/