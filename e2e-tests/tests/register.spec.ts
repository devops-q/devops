import {cleanupUser, expect, getUserByName, test} from './fixtures';

test.beforeEach(async ({dbPool}) => {
    await cleanupUser(dbPool, 'Me');
});

test.describe('User Registration Tests', () => {
    test('Register user via GUI', async ({page, dbPool}) => {
        // Navigate to register page
        await page.goto('/register');

        // Fill in form fields
        await page.fill('input[name=username]', 'Me');
        await page.fill('input[name=email]', 'me@some.where');
        await page.fill('input[name=password]', 'secure123');
        await page.fill('input[name=password2]', 'secure123');

        // Submit the form
        await page.keyboard.press('Enter');

        // Wait for flash message and check its content
        const flashMessage = await page.locator('.flashes').first();
        await expect(flashMessage).toBeVisible();
        await expect(flashMessage).toHaveText("You were successfully registered and can login now");
    });

    test('Register user via GUI and check DB entry', async ({page, dbPool}) => {
        // Check user doesn't exist yet
        let user = await getUserByName(dbPool, 'Me');
        expect(user).toBeNull();

        // Navigate to register page
        await page.goto('/register');

        // Fill in form fields
        await page.fill('input[name=username]', 'Me');
        await page.fill('input[name=email]', 'me@some.where');
        await page.fill('input[name=password]', 'secure123');
        await page.fill('input[name=password2]', 'secure123');

        // Submit the form
        await page.keyboard.press('Enter');

        // Wait for flash message and check its content
        const flashMessage = await page.locator('.flashes').first();
        await expect(flashMessage).toBeVisible();
        await expect(flashMessage).toHaveText("You were successfully registered and can login now");

        // Verify user exists in database
        user = await getUserByName(dbPool, 'Me');
        expect(user).not.toBeNull();
        expect(user.username).toBe('Me');
    });
});