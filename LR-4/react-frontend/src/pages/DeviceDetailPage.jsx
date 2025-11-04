import React from 'react';
import { Container, Button } from 'react-bootstrap';
import { useParams, useNavigate } from 'react-router-dom';

function DeviceDetailPage() {
  const { id } = useParams();
  const navigate = useNavigate();
  
  const handleBackClick = () => {
    navigate('/devices');
  };

  return (
    <Container>
      <Button variant="outline-primary" onClick={handleBackClick} className="mb-3">
        ← Назад к списку устройств
      </Button>
      
      <h1>Умное устройство #{id}</h1>
      <p className="text-muted">Детальная информация об устройстве</p>
      
      <div className="alert alert-info">
        Детальная страница устройства будет реализована на следующем этапе
      </div>
    </Container>
  );
}

export default DeviceDetailPage;