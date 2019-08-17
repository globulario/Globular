module.exports = {
  mode: "development",
  entry: "./services.js",
  output: {
    filename: 'services.js',
    library: 'globular',
    libraryTarget: 'window',
    libraryExport: 'default'
  }
};
