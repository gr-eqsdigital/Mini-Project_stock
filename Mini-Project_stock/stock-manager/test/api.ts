// using this import method, we don't need to have "const axios = require("axios");"
// at the test.ts file, axios does it with the import

import axios from "axios";

const api = axios.create({
    baseURL: "http://localhost:3000", // post
});

export default api;