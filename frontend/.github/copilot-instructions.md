<!-- Use this file to provide workspace-specific custom instructions to Copilot. For more details, visit https://code.visualstudio.com/docs/copilot/copilot-customization#_use-a-githubcopilotinstructionsmd-file -->

# Web Crawler Dashboard Project

This is a full-stack web crawler application with React TypeScript frontend and Go backend.

## Project Structure
- **Frontend**: React TypeScript with Vite, TailwindCSS, and modern UI components
- **Backend**: Go with Fiber framework, MySQL database, and web crawling capabilities

## Frontend Guidelines
- Use TypeScript for all code
- Implement responsive design with TailwindCSS
- Use React hooks and functional components
- Implement proper error handling and loading states
- Use socket.io-client for real-time updates
- Follow React best practices for state management

## Backend Guidelines
- Use Go with Fiber framework
- Implement proper error handling
- Use structured logging with Zap
- Follow clean architecture patterns (handler -> service -> repository)
- Implement web crawling with concurrent processing
- Use WebSockets for real-time communication

## Key Features to Implement
1. URL management (add, edit, delete)
2. Web crawling with analysis (HTML version, headings, links, forms)
3. Real-time progress tracking
4. Dashboard with charts and tables
5. Search and filtering capabilities
6. Bulk operations

## UI Components
- Use Headless UI for accessible components
- Use Heroicons for icons
- Use Recharts for data visualization
- Implement responsive layouts
- Follow modern design patterns
