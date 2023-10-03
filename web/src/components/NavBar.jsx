import { Navbar, Nav, NavDropdown } from 'react-bootstrap';
import { Link } from 'react-router-dom';

export default function MyNavbar() {
  return (
    <Navbar bg="light" expand="lg">
      <Navbar.Brand href="#home">Blogger</Navbar.Brand>
      <Navbar.Toggle aria-controls="basic-navbar-nav" />
      <Navbar.Collapse id="basic-navbar-nav">
        <Nav className="mr-auto">
          <Nav.Link href="/">Home</Nav.Link>
          <NavDropdown title="Code" id="basic-nav-dropdown">
            <NavDropdown.Item href="#action/1">Article 1</NavDropdown.Item>
            <NavDropdown.Item href="#action/2">Article 2</NavDropdown.Item>
          </NavDropdown>
          <Nav.Link
            href="https://github.com/BryceWayne/blogger"
            target="_blank"
            rel="noopener noreferrer">
              GitHub
          </Nav.Link>
        </Nav>
        <Nav className="ml-auto">
          <Nav.Link as={Link} to="/contact">Contact</Nav.Link>
        </Nav>
      </Navbar.Collapse>
    </Navbar>
  );
}
