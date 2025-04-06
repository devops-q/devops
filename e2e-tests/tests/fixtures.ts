import { test as base } from '@playwright/test';
import { hashSync } from 'bcrypt'; // Create a custom fixture
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

export async function createUser(pool: Pool, user: Partial<User>): Promise<User> {
  const hashedPassword = hashSync(user.pw_hash || 'default-hash', 10);
  const userId = user.id || Math.floor(Math.random() * 1000000) + 1;

  const result = await pool.query<User>(
    'insert into users (id, created_at, updated_at, deleted_at, username, email, pw_hash) values ($1, $2, $3, $4, $5, $6, $7) returning *',
    [
      user.id || userId,
      user.created_at || new Date(),
      user.updated_at || new Date(),
      user.deleted_at || null,
      user.username || 'default-username',
      user.email || 'default-email@example.com',
      hashedPassword,
    ],
  );

  return result.rows[0];
}

export async function createMessage(
  pool: Pool,
  userId: number,
  message: {
    text: string;
    flagged: boolean;
  },
): Promise<{ id: number; text: string; flagged: boolean; created_at: Date; updated_at: Date }> {
  const messageId = Math.floor(Math.random() * 1000000) + 1;

  const r = await pool.query(
    'insert into messages (id, created_at, updated_at, deleted_at, author_id, text, flagged) values ($1, $2, $3, $4, $5, $6, $7) returning *',
    [messageId, new Date(), new Date(), null, userId, message.text, message.flagged || false],
  );

  return r.rows[0];
}

export async function cleanupDb(pool: Pool): Promise<void> {
  await pool.query('DELETE FROM messages');
  await pool.query('DELETE FROM follower');
  await pool.query('DELETE FROM users');
  await pool.query('DELETE FROM api_users');
}
