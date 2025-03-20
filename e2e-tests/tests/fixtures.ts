import { test as base } from '@playwright/test';
import { Pool } from 'pg';
import { DB_URL } from '../config';

// Create a custom fixture
type TestFixtures = {
  dbPool: Pool;
};

type User = {
  id: string;
  created_at: Date;
  updated_at: Date;
  deleted_at: Date | null;
  username: string;
  email: string;
  pw_hash: string;
};

// Extend base test with our fixtures
export const test = base.extend<TestFixtures>({
  dbPool: async ({}, use) => {
    // Set up the database connection
    const pool = new Pool({
      connectionString: DB_URL,
    });

    // Verify connection is working
    await pool.query('SELECT NOW()');

    // Use the fixture
    await use(pool);

    // Clean up after tests
    await pool.end();
  },
});

export { expect } from '@playwright/test';

// Helper functions for database operations
export async function getUserByName(pool: Pool, name: string): Promise<User | null> {
  const result = await pool.query<User>('SELECT * FROM "users" WHERE username = $1', [name]);
  return result.rows[0] || null;
}

export async function cleanupUser(pool: Pool, name: string): Promise<void> {
  await pool.query('DELETE FROM "users" WHERE username = $1', [name]);
}
