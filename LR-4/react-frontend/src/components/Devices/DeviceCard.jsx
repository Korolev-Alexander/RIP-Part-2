import React from 'react';
import { Card, Button } from 'react-bootstrap';
import { useNavigate } from 'react-router-dom';

function DeviceCard({ device }) {
  const navigate = useNavigate();

  const handleDetailsClick = () => {
    navigate(`/devices/${device.id}`);
  };

  // Изображение по умолчанию если URL пустой
  const imageUrl = device.namespace_url || '/default-device.png';

  return (
    <Card className="h-100">
      <Card.Img 
        variant="top" 
        src={imageUrl}
        style={{ height: '200px', objectFit: 'contain', padding: '1rem' }}
        onError={(e) => {
          // Если изображение не загружается - ставим заглушку
          e.target.src = 'data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMjAwIiBoZWlnaHQ9IjIwMCIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj48cmVjdCB3aWR0aD0iMTAwJSIgaGVpZ2h0PSIxMDAlIiBmaWxsPSIjZGRkIi8+PHRleHQgeD0iNTAlIiB5PSI1MCUiIGZvbnQtc2l6ZT0iMTgiIGZpbGw9IiM5OTkiIGRvbWluYW50LWJhc2VsaW5lPSJtaWRkbGUiIHRleHQtYW5jaG9yPSJtaWRkbGUiPuKEoiDihpAg4oSWPC90ZXh0Pjwvc3ZnPg==';
        }}
      />
      
      <Card.Body className="d-flex flex-column">
        <Card.Title>{device.name}</Card.Title>
        <Card.Text className="flex-grow-1">
          {device.description}
        </Card.Text>
        
        <div className="mt-auto">
          <div className="mb-2">
            <small className="text-muted">
              <strong>Протокол:</strong> {device.protocol}
            </small>
            <br />
            <small className="text-muted">
              <strong>Трафик:</strong> {device.data_per_hour} Кб/ч
            </small>
          </div>
          
          <Button 
            variant="primary" 
            onClick={handleDetailsClick}
            className="w-100"
          >
            Подробнее
          </Button>
        </div>
      </Card.Body>
    </Card>
  );
}

// ДОБАВЛЯЕМ ЭТУ СТРОЧКУ!
export default DeviceCard;
