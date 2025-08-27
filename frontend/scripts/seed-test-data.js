// Script to seed test data for E2E tests
import { E2E_CONFIG } from '../tests/e2e.config.js';

async function seedTestData() {
  try {
    // First, login as admin to get token
    const loginResponse = await fetch(`${E2E_CONFIG.BACKEND_URL}/api/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        email: E2E_CONFIG.ADMIN_CREDENTIALS.email,
        password: E2E_CONFIG.ADMIN_CREDENTIALS.password,
      }),
    });

    if (!loginResponse.ok) {
      throw new Error('Failed to login as admin');
    }

    const loginData = await loginResponse.json();
    const adminToken = loginData.token;

    // Seed test users if needed
    const testUsers = [
      {
        name: 'Test User 1',
        email: 'testuser1@example.com',
        password: 'password123',
        role: 'user'
      },
      {
        name: 'Test User 2', 
        email: 'testuser2@example.com',
        password: 'password123',
        role: 'user'
      }
    ];

    for (const user of testUsers) {
      const userResponse = await fetch(`${E2E_CONFIG.BACKEND_URL}/api/debug/users/create`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${adminToken}`
        },
        body: JSON.stringify(user)
      });

      if (!userResponse.ok) {
        console.log(`User ${user.email} may already exist or failed to create`);
      }
    }

    console.log('Test data seeded successfully');
  } catch (error) {
    console.error('Failed to seed test data:', error);
  }
}

// Run if called directly
if (import.meta.url === `file://${process.argv[1]}`) {
  seedTestData();
}

export { seedTestData };