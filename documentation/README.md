# gitresolve Documentation Website

This folder contains the production documentation website for the gitresolve project.

The site documents:

- architecture and package boundaries
- end-to-end merge lifecycle
- command reference with real implementation behavior
- frontend and deployment workflow

## Tech Stack

- Next.js App Router
- React
- Tailwind CSS (v4 import style)

## Local Development

```bash
cd documentation
npm install
npm run dev
```

Default local URL: `http://localhost:3000`

## Quality and Build Commands

```bash
npm run lint     # static checks
npm run build    # production build output
npm run start    # run built app
```

## Frontend Structure

```text
documentation/
	src/
		app/
			globals.css   # design tokens + layout styles
			layout.js     # fonts + metadata + app shell
			page.js       # full documentation content
			robots.js     # crawler policy
			sitemap.js    # sitemap generation
```

## Production Notes

- Metadata is configured in `src/app/layout.js` for SEO and social previews.
- `robots.js` and `sitemap.js` are included for search indexing.
- The docs are intentionally source-of-truth aligned with current command implementations.
