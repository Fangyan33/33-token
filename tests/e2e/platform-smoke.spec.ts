import { expect, test } from "@playwright/test";

test("platform smoke route placeholders", async ({ page }) => {
  await page.goto("/");
  await expect(page).toHaveTitle(/Model API Platform/);
  await expect(page.getByText("低价模型 API 平台骨架")).toBeVisible();
});
