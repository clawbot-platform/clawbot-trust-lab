import { defineConfig } from "@playwright/test";

const port = 4173;
const baseURL = process.env.PLAYWRIGHT_BASE_URL || `http://127.0.0.1:${port}`;

export default defineConfig({
    testDir: "./tests/e2e",
    fullyParallel: false,
    retries: 0,
    use: {
        baseURL,
        headless: true,
        channel: process.env.PLAYWRIGHT_CHANNEL || "chrome",
        trace: "on-first-retry"
    },
    webServer: process.env.PLAYWRIGHT_BASE_URL
        ? undefined
        : {
            command: `npm run build && npm run preview -- --host 127.0.0.1 --port ${port}`,
            port,
            reuseExistingServer: !process.env.CI,
            timeout: 60_000
        }
});