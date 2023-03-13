import api from "./api";

let config;

// garantee auth token to test suits
beforeEach(async () => {
    const credentials = {
        "Email":"jtester@test.com",
		"Password":"123456789"
    };
	const response = await api.post("/login", credentials);
	// cookies = response.headers['set-cookie'];
    config = {
        headers:{
          Cookie: response.headers['set-cookie']
        }
      };
});

function randomString(): string {
    let r: string = (Math.random() + 1).toString(36).substring(7);
    // console.log("RANDOM: ", r);
    
    return r;
}

// user test suit
describe('Manage Users', () => {

    test('should succeed in creating a new user', async () => {
        expect.assertions(3);
        
        let user: any = {
            "Fname":"Mike",
            "Lname":"Tester",
            "Phone":"51515115151",
            // need to change the email, can't repeat email
            "Email": randomString() + "@test.com",
            "Password": "123456789"
        };
        
        // post request with new user information
        const response = await api.post("/signup", user);

        // test expectations
        expect(response.data).not.toBeNull();
        expect(response.data).toBeGreaterThanOrEqual(0);
        expect(response.status).toBe(200);
    });

    test('should succeed in calling the new user', async () => {
        expect.assertions(4);

        let user: any = {
            "fname":"Mike",
            "lname":"Tester",
            "phone":"51515115151",
            // allways needs a new email.
            "email": randomString() + "@test.com",
            "password": "123456789"
        };
        
        // post request with new user information
        const responseCreate = await api.post("/signup", user);
        expect(responseCreate.status).toBe(200);

        // Set user ID
        user.id = responseCreate.data;
        delete user.password;
        
        // get request with user email in url
        const response = await api.get("/user/get/" + user["email"], config);
        // const response = await api.get("/user/get/acoisa@test.com", config);
        // use data to extract the response info

        expect(response).not.toBeNull();
        expect(response.status).toBe(200);
        expect(response.data.message).toStrictEqual(user);
        
    })

    test('should succeed in deleting the test user', async () => {
        // post request to delete user
        expect.assertions(3);

        let user: any = {
            "fname":"Mike",
            "lname":"Tester",
            "phone":"51515115151",
            // allways needs a new email.
            "email": randomString() + "@test.com",
            "password": "123456789"
        };
        
        // post request with new user information
        const responseCreate = await api.post("/signup", user, config);
        expect(responseCreate.status).toBe(200);

        // set user ID
        user.id = responseCreate.data;
        
        // delete user
        // config object is used to ensure that the delete call, has the required credentials
        const responseDelete = await api.post("/user/delete", {
            "ID": user.id
        }, config);
        expect(responseDelete.status).toBe(200);
        expect(responseDelete).not.toBeNull();
    })
});