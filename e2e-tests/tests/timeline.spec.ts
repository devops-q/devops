import { describe } from 'node:test';
import { cleanupDb, createMessage, createUser, expect, test } from './fixtures';

test.afterAll(async ({ dbPool }) => {
  await cleanupDb(dbPool);
});

test.describe('Timeline Tests', () => {
  test('Public Timeline page loads', async ({ page }) => {
    // Navigate to timeline page
    await page.goto('/public');

    // Check if the timeline page is loaded
    const title = await page.title();
    expect(title).toBe('Public Timeline | MiniTwit');
  });

  test('Message is visible on timeline', async ({ page, dbPool }) => {
    const user = await createUser(dbPool, { username: 'user1', email: 'user1@example.com' });
    const message = await createMessage(dbPool, Number(user.id), {
      text: 'Hello, world!',
      flagged: false,
    });

    // Navigate to timeline page
    await page.goto('/public');

    // Check if the message is visible
    const messageElement = page.locator(`li#message-${message.id}`);
    await expect(messageElement).toBeVisible();
  });

  describe('When user locale is set to ja-JP and timezone is set to Asia/Tokyo', () => {
    test.use({
      locale: 'ja-JP',
      timezoneId: 'Asia/Tokyo',
    });

    test('Message createdAt is visible on timeline and is presented in the users Asia/Tokyo timezone and ja-JP locale', async ({
      page,
      dbPool,
    }) => {
      const user = await createUser(dbPool, { username: 'user2', email: 'user2@example.com' });
      const message = await createMessage(dbPool, Number(user.id), {
        text: 'Hello, world! What time is it in my timezone?',
        flagged: false,
      });

      // Navigate to timeline page
      await page.goto('/public');
      //   find the corresponding li element
      const messageElement = page.locator(`li#message-${message.id}`);
      const timeElement = messageElement.locator('small.time');
      const timeText = await timeElement.innerText();
      // Check if the createdAt time is formatted correctly
      const date = message.created_at;
      const formattedDate = `Written @ ${date.toLocaleTimeString('ja-JP', {
        timeZone: 'Asia/Tokyo',
      })}, ${date.toLocaleDateString('ja-JP', {
        timeZone: 'Asia/Tokyo',
      })}`;

      expect(timeText).toContain(formattedDate);
    });
  });
});
describe('When user locale is set to da-DK and timezone is set to Europe/Copenhagen', () => {
  test.use({
    locale: 'da-DK',
    timezoneId: 'Europe/Copenhagen',
  });

  test('Message createdAt is visible on timeline and is presented in the Europe/Copenhagen timezone and da-DK locale', async ({
    page,
    dbPool,
  }) => {
    const user = await createUser(dbPool, { username: 'user2', email: 'user2@example.com' });
    const message = await createMessage(dbPool, Number(user.id), {
      text: 'Hello, world! What time is it in my timezone?',
      flagged: false,
    });

    // Navigate to timeline page
    await page.goto('/public');
    //   find the corresponding li element
    const messageElement = page.locator(`li#message-${message.id}`);
    const timeElement = messageElement.locator('small.time');
    const timeText = await timeElement.innerText();
    // Check if the createdAt time is formatted correctly
    const date = message.created_at;
    const formattedDate = `Written @ ${date.toLocaleTimeString('da-DK', {
      timeZone: 'Europe/Copenhagen',
    })}, ${date.toLocaleDateString('da-DK', {
      timeZone: 'Europe/Copenhagen',
    })}`;

    expect(timeText).toContain(formattedDate);
  });
});
