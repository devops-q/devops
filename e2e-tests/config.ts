import * as path from 'node:path';

export const HEALTH_CHECK_URL = 'http://localhost:80/api/v1/health';
export const GUI_URL = 'http://localhost:80';
export const HEALTH_CHECK_TIMEOUT = 120000; // 2 minutes in milliseconds
export const HEALTH_CHECK_INTERVAL = 5000; // 5 seconds between attempts
export const DOCKER_COMPOSE_PATH = path.resolve(__dirname, './docker-compose.yml');
export const DB_URL = 'postgresql://postgres:postgres@localhost:5432/test';
