var path = require('path')
var webpack = require('webpack')

module.exports = {
  entry: [
    "./js/app.jsx"
  ],
  output: {
    path: __dirname + '/public',
    filename: "bundle.js"
  },
  module: {
    rules: [
      {
        test: /.jsx?$/,
        exclude: /node_modules/,
        use: {
          loader: 'babel-loader',
          options: {
            presets: ['es2015', 'react']
          }
        }
      },
      { test: /\.js$/, exclude: /node_modules/, loader: 'babel-loader'},
      { test: /\.css$/, loader: "style!css" }
    ]
  },
  plugins: [
  ],
  devServer: {
    contentBase:  path.join(__dirname, "public"),
    proxy: {
      '/api': {
        //target: 'http://localhost:8081',
        target: 'https://particle-volt.appspot.com',
        secure: false,
        changeOrigin: true
      }
    }
  }
};
