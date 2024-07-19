exports.config = {
    tests: './*.test.js',
    output: './output',
    helpers: {
      Puppeteer: {
        url: 'http://localhost:3000',
        show: true,
        windowSize: '1200x900'
      }
    },
    include: {
      I: './steps_file.js'
    },
    bootstrap: null,
    mocha: {},
    name: 'frontend'
  };
  