import MyNavbar from './components/NavBar';
import Home from './pages/Home';
import AboutUs from './pages/AboutUs';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';

function App() {
  return (
    <Router>
      <MyNavbar />
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/contact" element={<AboutUs />} />
      </Routes>
    </Router>
  );
}

export default App;
