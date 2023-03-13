"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
var __generator = (this && this.__generator) || function (thisArg, body) {
    var _ = { label: 0, sent: function() { if (t[0] & 1) throw t[1]; return t[1]; }, trys: [], ops: [] }, f, y, t, g;
    return g = { next: verb(0), "throw": verb(1), "return": verb(2) }, typeof Symbol === "function" && (g[Symbol.iterator] = function() { return this; }), g;
    function verb(n) { return function (v) { return step([n, v]); }; }
    function step(op) {
        if (f) throw new TypeError("Generator is already executing.");
        while (g && (g = 0, op[0] && (_ = 0)), _) try {
            if (f = 1, y && (t = op[0] & 2 ? y["return"] : op[0] ? y["throw"] || ((t = y["return"]) && t.call(y), 0) : y.next) && !(t = t.call(y, op[1])).done) return t;
            if (y = 0, t) op = [op[0] & 2, t.value];
            switch (op[0]) {
                case 0: case 1: t = op; break;
                case 4: _.label++; return { value: op[1], done: false };
                case 5: _.label++; y = op[1]; op = [0]; continue;
                case 7: op = _.ops.pop(); _.trys.pop(); continue;
                default:
                    if (!(t = _.trys, t = t.length > 0 && t[t.length - 1]) && (op[0] === 6 || op[0] === 2)) { _ = 0; continue; }
                    if (op[0] === 3 && (!t || (op[1] > t[0] && op[1] < t[3]))) { _.label = op[1]; break; }
                    if (op[0] === 6 && _.label < t[1]) { _.label = t[1]; t = op; break; }
                    if (t && _.label < t[2]) { _.label = t[2]; _.ops.push(op); break; }
                    if (t[2]) _.ops.pop();
                    _.trys.pop(); continue;
            }
            op = body.call(thisArg, _);
        } catch (e) { op = [6, e]; y = 0; } finally { f = t = 0; }
        if (op[0] & 5) throw op[1]; return { value: op[0] ? op[1] : void 0, done: true };
    }
};
exports.__esModule = true;
var api_1 = require("./api");
var config;
// garantee auth token to test suits
beforeEach(function () { return __awaiter(void 0, void 0, void 0, function () {
    var credentials, response;
    return __generator(this, function (_a) {
        switch (_a.label) {
            case 0:
                credentials = {
                    "Email": "jtester@test.com",
                    "Password": "123456789"
                };
                return [4 /*yield*/, api_1["default"].post("/login", credentials)];
            case 1:
                response = _a.sent();
                // cookies = response.headers['set-cookie'];
                config = {
                    headers: {
                        Cookie: response.headers['set-cookie']
                    }
                };
                return [2 /*return*/];
        }
    });
}); });
function randomString() {
    var r = (Math.random() + 1).toString(36).substring(7);
    // console.log("RANDOM: ", r);
    return r;
}
// user test suit
describe('Manage Users', function () {
    test('should succeed in creating a new user', function () { return __awaiter(void 0, void 0, void 0, function () {
        var user, response;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0:
                    expect.assertions(3);
                    user = {
                        "Fname": "Mike",
                        "Lname": "Tester",
                        "Phone": "51515115151",
                        // need to change the email, can't repeat email
                        "Email": randomString() + "@test.com",
                        "Password": "123456789"
                    };
                    return [4 /*yield*/, api_1["default"].post("/signup", user)];
                case 1:
                    response = _a.sent();
                    // test expectations
                    expect(response.data).not.toBeNull();
                    expect(response.data).toBeGreaterThanOrEqual(0);
                    expect(response.status).toBe(200);
                    return [2 /*return*/];
            }
        });
    }); });
    test('should succeed in calling the new user', function () { return __awaiter(void 0, void 0, void 0, function () {
        var user, responseCreate, response;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0:
                    expect.assertions(4);
                    user = {
                        "fname": "Mike",
                        "lname": "Tester",
                        "phone": "51515115151",
                        // allways needs a new email.
                        "email": randomString() + "@test.com",
                        "password": "123456789"
                    };
                    return [4 /*yield*/, api_1["default"].post("/signup", user)];
                case 1:
                    responseCreate = _a.sent();
                    expect(responseCreate.status).toBe(200);
                    // Set user ID
                    user.id = responseCreate.data;
                    delete user.password;
                    return [4 /*yield*/, api_1["default"].get("/user/get/" + user["email"], config)];
                case 2:
                    response = _a.sent();
                    // const response = await api.get("/user/get/acoisa@test.com", config);
                    // use data to extract the response info
                    expect(response).not.toBeNull();
                    expect(response.status).toBe(200);
                    expect(response.data.message).toStrictEqual(user);
                    return [2 /*return*/];
            }
        });
    }); });
    test('should succeed in deleting the test user', function () { return __awaiter(void 0, void 0, void 0, function () {
        var user, responseCreate, responseDelete;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0:
                    // post request to delete user
                    expect.assertions(3);
                    user = {
                        "fname": "Mike",
                        "lname": "Tester",
                        "phone": "51515115151",
                        // allways needs a new email.
                        "email": randomString() + "@test.com",
                        "password": "123456789"
                    };
                    return [4 /*yield*/, api_1["default"].post("/signup", user, config)];
                case 1:
                    responseCreate = _a.sent();
                    expect(responseCreate.status).toBe(200);
                    // set user ID
                    user.id = responseCreate.data;
                    return [4 /*yield*/, api_1["default"].post("/user/delete", {
                            "ID": user.id
                        }, config)];
                case 2:
                    responseDelete = _a.sent();
                    expect(responseDelete.status).toBe(200);
                    expect(responseDelete).not.toBeNull();
                    return [2 /*return*/];
            }
        });
    }); });
});
