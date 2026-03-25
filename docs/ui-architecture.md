# UI Architecture

Phase 8 adds a thin operator UI under `web/`.

## Split of responsibilities

### Go backend

The backend continues to own:

- round execution
- promotion records
- detection results
- report artifact generation
- operator review persistence
- comparison logic

### React UI

The UI only owns:

- fetching operator APIs
- presenting rounds, promotions, and reports
- submitting review actions
- rendering existing Markdown and JSON artifacts

## Pages

The UI includes:

- Rounds
- Round Detail
- Promotions
- Promotion / Detection Explorer
- Reports

## Design intent

The UI is intentionally restrained:

- one sidebar navigation
- review-focused cards and tables
- simple status pills
- no chart-heavy dashboarding

## Local development

The UI uses Vite with a dev proxy to the Go backend on `http://127.0.0.1:8090`.

This keeps the frontend thin and avoids duplicating API or report logic in the browser.
