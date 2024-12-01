# Usage
### Backend
- Make sure you have `Go` installed
- Navigate to the folder where you cloned this repo and type the following in your terminal
  - `go mod init tech-interview-buddy`
  - `go mod tidy`
  - `go run cmd/api/main.go`
- Note: `go mod` commands are only required for the first use.

### Frontend
- Make sure you have `Node.js` installed
- From the folder you cloned this repo type the following in your terminal
  - `cd frontend/code-connect`
  - `npm install`
  - `npm run dev`
- Add a `.env.local` and `firebase-admin-config.json` in the frontend/code-connect directory to set the API keys for the firebase functionality
- Note: `npm install` is only required for the first use.

### Jitsi Server
- Our Jitsi Server is self-hosted using the following guide: [https://jitsi.github.io/handbook/docs/devops-guide/devops-guide-quickstart/](https://jitsi.github.io/handbook/docs/devops-guide/devops-guide-docker)
- Any references to "actual-terribly-longhorn.ngrok-free.app" are a reference to our ngrok tunnel running the self-hosted jitsi server

# Public API Endpoints
- Documentation can be found [here](https://github.com/rebeccaborissova/tech-interview-buddy/wiki/Public-API-Endpoints)

# Jira Board
[link!](https://rebeccaborissov.atlassian.net/jira/software/projects/SCRUM/boards/1/timeline)
