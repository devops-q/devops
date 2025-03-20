import {test as teardown} from '@playwright/test';
import {exec} from 'child_process';
import {promisify} from 'util';
import {DOCKER_COMPOSE_PATH} from "../config";

const execAsync = promisify(exec);

teardown('Tear down Docker Compose environment', async () => {
    console.log('Tearing down Docker Compose setup...');

    try {
        await execAsync(`docker compose -f ${DOCKER_COMPOSE_PATH} down`);
        console.log('Docker Compose teardown completed successfully');
    } catch (error) {
        console.error('Failed to tear down Docker Compose setup:', error);
        throw error;
    }
});