#!/bin/bash
npm i swagger-ui-dist
mv ./node_modules/swagger-ui-dist ./static
rm -R ./node_modules
cd ./static
mv -t ../cmd/rest_server/static swagger-ui-bundle.js swagger-ui-standalone-preset.js swagger-initializer.js favicon-16x16.png favicon-32x32.png index.css swagger-ui.css index.html
cd ..
rm -R ./static

rm package-lock.json
rm package.json