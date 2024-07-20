exports.config = {
    tests: './*.test.js',
    output: './output',
    helpers: {
      Puppeteer: {
        url: 'https://panaderia-stock-frontend-app-6df615a13979.herokuapp.com',
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
  