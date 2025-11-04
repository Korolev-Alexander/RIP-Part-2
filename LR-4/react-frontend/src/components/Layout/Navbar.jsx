import React from 'react';
import { Navbar, Nav, Container } from 'react-bootstrap';
import { LinkContainer } from 'react-router-bootstrap';

function AppNavbar() {
  return (
    <Navbar bg="primary" variant="dark" expand="lg">
      <Container>
        <Navbar.Brand href="/">Система Умный Дом</Navbar.Brand>
        <Nav className="me-auto">
          <LinkContainer to="/">
            <Nav.Link>Главная</Nav.Link>
          </LinkContainer>
          <LinkContainer to="/devices">
            <Nav.Link>Умные устройства</Nav.Link>
          </LinkContainer>
        </Nav>
      </Container>
    </Navbar>
  );
}

export default AppNavbar;