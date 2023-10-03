import '../index.css';

export default function Home() {
  return (
    <div className="d-flex flex-column align-items-center justify-content-center home-container" style={{ height: '100vh' }}>
      <h1 className="title-animate">Welcome to Blogger!</h1>
      <div className="circle-animate">
        <svg height="100" width="100">
          <circle cx="50" cy="50" r="40" stroke="black" strokeWidth="3" fill="red" />
        </svg>
      </div>
      <p>Your go-to platform for self-documenting APIs</p>
    </div>
  );
}
