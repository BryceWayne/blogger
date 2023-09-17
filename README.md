# Blogger - Code Blogger

## Project Overview

Automatically generates blog posts that reflect the changes and updates within the codebase. Utilizes chatGPT to convert GitHub push events into meaningful content that documents project progress.

## Features

- Uses OpenAI's GPT to generate content
- Automatically triggered by GitHub push events
- Real-time update and documentation of project progress
- Stored in Firestore for scalability and real-time updates

## Tech Stack

- GCP (Google Cloud Platform)
- Go Fiber for API
- React for Frontend
- Firestore as Database
- OpenAI API for Content Generation

## Installation

1. Clone the repository  
```bash
git clone [repository_link]
```

2. Deploy API  
```bash
./deploy_api.sh $PROJECT_ID
```

3. Deploy Web  (coming soon)
```bash
./deploy_web.sh $PROJECT_ID
```

## Usage

- Set up GCP Firestore for data storage.
- Connect your GitHub repository to webhooks.
- Create an OpenAI account and get an API key for content generation.
- Connect the system to your webpage to display blog posts.

## Contributions

Feel free to contribute to front-end development or work on scaling the application. Check out the `issues` tab for open tasks or bugs.

## License

No license specified as of yet.
