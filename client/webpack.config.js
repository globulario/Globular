module.exports = {
  mode: "development",
  entry: "./services.js",
  output: {
    filename: 'globularServices.js',
    library: 'globular',
    libraryTarget: 'window',
    libraryExport: 'default'
  }
};
