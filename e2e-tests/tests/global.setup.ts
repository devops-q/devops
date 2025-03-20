import { exec } from 'node:child_process';
import { promisify } from 'node:util';
import { test as setup } from '@playwright/test';
import {
  DOCKER_COMPOSE_PATH,
  HEALTH_CHECK_INTERVAL,
  HEALTH_CHECK_TIMEOUT,
  HEALTH_CHECK_URL,
} from '../config';

const execAsync = promisify(exec);

/**
 * Sleep function to wait between polling attempts
 */
const sleep = (ms: number) => new Promise((resolve) => setTimeout(resolve, ms));

/**
 * Poll the health endpoint until it returns 200 or timeout is reached
 */
async function waitForHealthEndpoint(
  url: string,
  timeout: number,
  interval: number,
): Promise<boolean> {
  const startTime = Date.now();
  let attemptCount = 0;

  while (Date.now() - startTime < timeout) {
    attemptCount++;
    try {
      const response = await fetch(url);
      if (response.status === 200) {
        console.log(`Health check successful after ${attemptCount} attempts`);
        return true;
      }
      console.log(`Health check attempt ${attemptCount}: received status ${response.status}`);
    } catch (error) {
      console.log(`Health check attempt ${attemptCount} failed: ${error.message}`);
    }

    await sleep(interval);
  }

  throw new Error(`Health check timed out after ${timeout}ms`);
}

setup('Start Docker Compose environment and wait for API health', async ({}, testInfo) => {
  testInfo.setTimeout(180000);
  try {
    // Start containers
    await execAsync(`docker compose -f ${DOCKER_COMPOSE_PATH} up -d`);
    console.log('Docker Compose containers started');

    // Wait for API to be ready by polling the health endpoint
    console.log(`Waiting for API health check at ${HEALTH_CHECK_URL}`);
    await waitForHealthEndpoint(HEALTH_CHECK_URL, HEALTH_CHECK_TIMEOUT, HEALTH_CHECK_INTERVAL);
    console.log('Application is ready for testing');
  } catch (error) {
    console.error('Failed during setup:', error);

    // Try to clean up on failure
    try {
      await execAsync(`docker compose -f ${DOCKER_COMPOSE_PATH} down`);
      console.log('Cleaned up Docker Compose environment after failure');
    } catch (cleanupError) {
      console.error('Failed to clean up Docker Compose after failure:', cleanupError);
    }

    throw error;
  }
});
