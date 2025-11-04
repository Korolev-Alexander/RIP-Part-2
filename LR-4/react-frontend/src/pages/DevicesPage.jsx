import React from 'react';
import { Container } from 'react-bootstrap';

function DevicesPage() {
  return (
    <Container>
      <h1>Каталог умных устройств</h1>
      <p className="text-muted">
        Выберите подходящие устройства для вашего умного дома
      </p>
      <div className="alert alert-info">
        Список устройств и фильтры будут добавлены на следующем этапе
      </div>
    </Container>
  );
}

export default DevicesPage;