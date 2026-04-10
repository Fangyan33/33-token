import { expect, test } from 'playwright/test';

test('platform smoke renders the shell', async ({ page }) => {
  await page.goto('/');

  await expect(page.getByTestId('platform-shell')).toBeVisible();
  await expect(page.getByRole('heading', { name: 'Platform MVP' })).toBeVisible();
  await expect(page.getByText('Local stack ready')).toBeVisible();
});
