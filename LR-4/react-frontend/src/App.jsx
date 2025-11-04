import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { Container } from 'react-bootstrap';
import 'bootstrap/dist/css/bootstrap.min.css';

import Navbar from './components/Layout/Navbar';
import HomePage from './pages/HomePage';
import DevicesPage from './pages/DevicesPage';
import DeviceDetailPage from './pages/DeviceDetailPage';

function App() {
  return (
    <Router>
      <Navbar />
      <Container fluid className="py-4">
        <Routes>
          <Route path="/" element={<HomePage />} />
          <Route path="/devices" element={<DevicesPage />} />
          <Route path="/devices/:id" element={<DeviceDetailPage />} />
        </Routes>
      </Container>
    </Router>
  );
}

export default App;