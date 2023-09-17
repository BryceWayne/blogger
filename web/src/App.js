import React from 'react';
import './App.css';


import { Light as SyntaxHighlighter } from 'react-syntax-highlighter';
import { docco } from 'react-syntax-highlighter/dist/esm/styles/hljs';

const BeautifyGoCode = ({ code }) => {
  const parsed = code.split('\n').map((line, idx) => {
    if (line.startsWith("//")) {
      return <div className="comment" key={idx}>{line}</div>;
    } else {
      return <SyntaxHighlighter language="go" style={docco} key={idx}>{line}</SyntaxHighlighter>;
    }
  });
  return <div>{parsed}</div>;
};

const Blog = ({ text }) => {
  const blocks = text.split("```");
  return (
    <div>
      {blocks.map((block, index) => {
        if (index % 2 === 0) {
          return <p key={index}>{block}</p>;
        } else {
          return <BeautifyGoCode key={index} code={block.trim()} />;
        }
      })}
    </div>
  );
};

const blogText = `
  This is my blog.
  Here's some Go code:
  \`\`\`
  // Adds two numbers
  func add(a int, b int) int {
    return a + b
  }
  \`\`\`
  Hope you like it.
`;

function App() {
  return (
    <div className="App">
      <header className="App-header">
        <h1>BlogGPT</h1>
      </header>
      
      <main className="App-main">
        <h2>Featured Articles</h2>
        <div className="article-card">
          <h3>Example Go</h3>
          <Blog text={blogText} />
          <p>The Go code defines a simple function, add, that takes two integers as arguments and returns their sum.</p>
        </div>
        <div className="article-card">
          <h3>Article 2</h3>
          <p>This is a summary of the second article...</p>
        </div>
      </main>
      
      <footer className="App-footer">
        <p>Rektor.tech Copyright 2023. All Rights Reserved.</p>
      </footer>
    </div>
  );
}

export default App;