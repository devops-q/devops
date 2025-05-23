import { defineConfig, devices } from '@playwright/test';
import { GUI_URL } from './config'; /**
 * Read environment variables from file.
 * https://github.com/motdotla/dotenv
 */

/**
 * Read environment variables from file.
 * https://github.com/motdotla/dotenv
 */
// import dotenv from 'dotenv';
// import path from 'path';
// dotenv.config({ path: path.resolve(__dirname, '.env') });

/**
 * See https://playwright.dev/docs/test-configuration.
 */
export default defineConfig({
  testDir: './tests',
  /* Run tests in files in parallel */
  fullyParallel: false,
  /* Fail the build on CI if you accidentally left test.only in the source code. */
  forbidOnly: !!process.env.CI,
  /* Retry on CI only */
  retries: process.env.CI ? 2 : 0,
  /* Opt out of parallel tests on CI. */
  workers: 1,
  /* Reporter to use. See https://playwright.dev/docs/test-reporters */
  reporter: 'html',
  /* Shared settings for all the projects below. See https://playwright.dev/docs/api/class-testoptions. */
  use: {
    /* Base URL to use in actions like `await page.goto('/')`. */
    baseURL: GUI_URL,

    /* Collect trace when retrying the failed test. See https://playwright.dev/docs/trace-viewer */
    trace: 'on-first-retry',
    headless: true,
    viewport: { width: 1280, height: 720 },
    ignoreHTTPSErrors: true,
    video: 'on-first-retry',
  },

  timeout: 30000,

  expect: {
    timeout: 5000,
  },

  /* Configure projects for major browsers */
  projects: [
    {
      name: 'setup',
      testMatch: /global\.setup\.ts/,
      teardown: 'teardown',
    },
    {
      name: 'teardown',
      testMatch: /global\.teardown\.ts/,
    },
    {
      name: 'chromium',
      use: {
        browserName: 'chromium',
      },
      dependencies: ['setup'],
    },

    /* Test against branded browsers. */
    {
      name: 'Google Chrome',
      use: { ...devices['Desktop Chrome'], channel: 'chrome' },
      dependencies: ['setup'],
    },
    {
      name: 'Safari',
      use: { ...devices['Desktop Safari'] },
      dependencies: ['setup'],
    },

    /* Test against mobile devices. */
    {
      name: 'iPhone 12',
      use: { ...devices['iPhone 12'] },
      dependencies: ['setup'],
    },
    {
      name: 'Pixel 5',
      use: { ...devices['Pixel 5'] },
      dependencies: ['setup'],
    },
    {
      name: 'Galaxy S20',
      use: { ...devices['Galaxy S20'] },
      dependencies: ['setup'],
    },
  ],

  /* Run your local dev server before starting the tests */
  // webServer: {
  //   command: 'npm run start',
  //   url: 'http://127.0.0.1:3000',
  //   reuseExistingServer: !process.env.CI,
  // },
});
