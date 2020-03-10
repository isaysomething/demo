const minify = require('@node-minify/core');
const gcc = require('@node-minify/google-closure-compiler');
const cleanCSS = require('@node-minify/clean-css');


minify({
    compressor: gcc,
    input: 'web/js/app.js',
    output: 'web/js/app.min.js',
    callback: function(err, min) {
        if (err) {
            console.error(err)
        }
    }
  });
  
minify({
    compressor: cleanCSS,
    input: 'web/css/app.css',
    output: 'web/css/app.min.css',
    callback: function(err, min) {
        if (err) {
            console.error(err)
        }
    }
  });
  