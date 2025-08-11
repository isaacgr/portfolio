//
// update-template.js
// Credit: https://github.com/tailwindlabs/tailwindcss/discussions/2743

const fs = require("fs");

const templateFile = "./web/views/base.html";
const controlFile = "./web/static/css/output.css";

const version = fs.statSync(controlFile).mtimeMs;
const replacement = `href="/static/css/output.css?v=${version}"`;

const searchRegex = /href=\"\/static\/css\/output\.css.*\"/g;

const htmlTemplate = fs
  .readFileSync(templateFile, "utf-8")
  .replace(searchRegex, replacement);

fs.writeFileSync(templateFile, htmlTemplate, {
  encoding: "utf-8"
});
