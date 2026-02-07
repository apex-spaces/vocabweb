# VocabWeb Frontend

A modern vocabulary learning platform built with Next.js 14, TypeScript, and Tailwind CSS.

## Tech Stack

- **Framework**: Next.js 14 (App Router)
- **Language**: TypeScript
- **Styling**: Tailwind CSS
- **UI Components**: shadcn/ui (to be added)
- **Fonts**: Playfair Display + DM Sans (Google Fonts)

## Getting Started

### Prerequisites

- Node.js 18+ 
- npm or yarn

### Installation

```bash
# Install dependencies
npm install

# Copy environment variables
cp .env.example .env

# Run development server
npm run dev
```

Open [http://localhost:3000](http://localhost:3000) to view the app.

## Project Structure

```
frontend/
├── src/
│   ├── app/              # Next.js App Router pages
│   ├── components/       # React components
│   │   ├── ui/          # shadcn/ui components
│   │   ├── layout/      # Layout components (Sidebar, Header)
│   │   └── common/      # Common components (Logo)
│   ├── lib/             # Utilities and API client
│   └── types/           # TypeScript type definitions
├── public/              # Static assets
└── Dockerfile           # Multi-stage Docker build
```

## Available Scripts

- `npm run dev` - Start development server
- `npm run build` - Build for production
- `npm start` - Start production server
- `npm run lint` - Run ESLint

## Deployment

### Docker

```bash
docker build -t vocabweb-frontend .
docker run -p 3000:3000 vocabweb-frontend
```

### Google Cloud Run

```bash
gcloud run deploy vocabweb-frontend \
  --source . \
  --region asia-east2 \
  --platform managed
```

## Design System

- **Background**: #0F172A (slate-900)
- **Card**: #1E293B (slate-800)
- **Accent**: #F59E0B (amber-500)
- **Primary Font**: Playfair Display (serif)
- **Body Font**: DM Sans (sans-serif)

## Features (Planned)

- ✅ Landing page
- ✅ Authentication (placeholder)
- ✅ Dashboard with sidebar navigation
- ✅ 12 main pages (placeholders)
- ⏳ Word management
- ⏳ Spaced repetition review
- ⏳ Progress tracking
- ⏳ Collections and categories
- ⏳ Learning schedule
- ⏳ Achievements system
- ⏳ Community features

## License

MIT
