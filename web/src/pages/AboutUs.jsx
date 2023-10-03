export default function AboutUs() {
  return (
    <div className="d-flex flex-column align-items-center justify-content-center" style={{ height: '100vh' }}>
      <h1>About Blogger</h1>
      <p>A Self-Documenting API Platform</p>
      
      <div className="text-left">
        <h4>Core Features:</h4>
        <ul>
          <li>Automated API Documentation</li>
          <li>GitHub Webhook Integration</li>
          <li>AI-powered by OpenAI's ChatGPT</li>
        </ul>
        
        <h4>Technology Stack:</h4>
        <p>Go, Go-Fiber, GitHub Webhooks, OpenAI</p>
        
        <h4>Availability:</h4>
        <p>Free*, Open Source, and ready to be deployed.</p>
      </div>
      
      <small>*Free with conditions. Please refer to the license.</small>
    </div>
  );
}
