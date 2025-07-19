# QKO Marketing Site

## Privacy Protection Setup

This site is currently protected from search engine indexing and AI crawlers until ready for launch.

### Protection Layers:
- `robots.txt` - Blocks all web crawlers
- Meta robots tags - Prevents indexing and following
- HTTP headers - Server-level protection
- Cache prevention - Prevents content caching

### To Launch:
1. Remove or modify `public/robots.txt`
2. Remove meta robots tags from `src/layouts/Layout.astro`
3. Remove server headers from `astro.config.mjs`
4. Remove `.htaccess` restrictions
5. Submit to search engines when ready

### Files to Track:
- ✅ `public/robots.txt` - Standard web config
- ✅ `public/.htaccess` - Apache server config
- ✅ `src/layouts/Layout.astro` - Meta robots tags
- ✅ `astro.config.mjs` - Server headers

### Files NOT to Track:
- ❌ `.htpasswd` - Password files (if added)
- ❌ `.env` - Environment secrets
- ❌ `secrets.json` - API keys, etc.

## Development

```bash
npm install
npm run dev
npm run build
npm run preview
```
