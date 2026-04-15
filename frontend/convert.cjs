const sharp = require('sharp');
const fs = require('fs');

sharp('./src/assets/logo.svg')
  .resize(256, 256)
  .png()
  .toFile('../build/appicon.png')
  .then(() => console.log('Successfully created appicon.png for Wails'))
  .catch(err => console.error('Error generating icon:', err));
