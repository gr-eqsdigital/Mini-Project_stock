import api from "./api";

// login test suit
describe('Login', () => {

    test('should succeed with valid credentials', async () => {
        expect.assertions(3);
        // post request with valid credentials
            const credentials = {
            "Email":"jtester@test.com",
            "Password":"123456789"
            };
        const response = await api.post("/login", credentials);
        // set cookies, to check auth token
        const cookies = response.headers['set-cookie'];
        // test assertions
        
        // criar mais verificações para validar cookies
        // authorization, if it changes tests would pass
        expect(cookies).not.toBeNull();
        expect(response).not.toBeNull();
        expect(response.status).toBe(200);
    });

    test('should fail with invalid credentials', async () => {
		 expect.assertions(1);

		try {
            // post request with invalid credentials
            const response = await api.post("/login", {
                "Email":"invalid",
                "Password":"invalid"
            });
            expect(response).not.toBeNull();
        } catch (error) {
            // assertion for error response
            expect(error.response.status).toBe(401);
        }
    });
});