"use strict";
// using this import method, we don't need to have "const axios = require("axios");"
// at the test.ts file, axios does it with the import
exports.__esModule = true;
var axios_1 = require("axios");
var api = axios_1["default"].create({
    baseURL: "http://localhost:3000"
});
exports["default"] = api;
