# Frontend README

This directory contains the Vue 3 frontend application for ZiXiao Git Server.

## Quick Start

```bash
# Install dependencies
npm install

# Development server (with hot reload)
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview
```

## Tech Stack

- **Vue 3** - Progressive JavaScript framework
- **Vue Router** - Official router
- **Pinia** - State management
- **Axios** - HTTP client
- **Vite** - Build tool
- **MDUI 2.x** - Material Design UI

## Project Structure

```
src/
├── views/          # Page components
├── stores/         # Pinia stores
├── services/       # API services
├── router/         # Vue Router config
├── utils/          # Utilities
├── App.vue         # Root component
└── main.js         # Entry point
```

## Documentation

- [Frontend Development Guide](../docs/FRONTEND_DEV.md)
- [Deployment Guide](../docs/FRONTEND_DEPLOYMENT.md)

## Scripts

- `npm run dev` - Start development server on port 3000
- `npm run build` - Build for production (outputs to `../web/dist`)
- `npm run preview` - Preview production build locally

## Environment Variables

Create `.env.local` for local overrides:

```env
VITE_API_BASE_URL=http://localhost:8080/api/v1
```

## License

MIT
