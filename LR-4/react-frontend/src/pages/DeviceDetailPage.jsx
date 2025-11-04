import React, { useState, useEffect } from 'react';
import { Container, Row, Col, Card, Button, Alert } from 'react-bootstrap';
import { useParams, useNavigate } from 'react-router-dom';
import Breadcrumbs from '../components/Layout/Breadcrumbs';
import { apiService } from '../services/api';

function DeviceDetailPage() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [device, setDevice] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const breadcrumbsItems = [
    { label: 'Умные устройства', href: '/devices' },
    { label: device?.name || 'Загрузка...', active: true }
  ];

  useEffect(() => {
    loadDevice();
  }, [id]);

  const loadDevice = async () => {
    setLoading(true);
    setError(null);
    try {
      const deviceData = await apiService.getDevice(id);
      if (deviceData) {
        setDevice(deviceData);
      } else {
        setError('Устройство не найдено');
      }
    } catch (error) {
      setError('Ошибка при загрузке устройства');
      console.error('Error loading device:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleBackClick = () => {
    navigate('/devices');
  };

  if (loading) {
    return (
      <Container>
        <Breadcrumbs items={breadcrumbsItems} />
        <div className="text-center py-5">
          <div className="spinner-border" role="status">
            <span className="visually-hidden">Загрузка...</span>
          </div>
        </div>
      </Container>
    );
  }

  if (error || !device) {
    return (
      <Container>
        <Breadcrumbs items={breadcrumbsItems} />
        <Alert variant="danger">
          {error || 'Устройство не найдено'}
        </Alert>
        <Button variant="primary" onClick={handleBackClick}>
          Вернуться к списку устройств
        </Button>
      </Container>
    );
  }

  const imageUrl = device.namespace_url || 'data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iNDAwIiBoZWlnaHQ9IjQwMCIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj48cmVjdCB3aWR0aD0iMTAwJSIgaGVpZ2h0PSIxMDAlIiBmaWxsPSIjZGRkIi8+PHRleHQgeD0iNTAlIiB5PSI1MCUiIGZvbnQtc2l6ZT0iMjQiIGZpbGw9IiM5OTkiIGRvbWluYW50LWJhc2VsaW5lPSJtaWRkbGUiIHRleHQtYW5jaG9yPSJtaWRkbGUiPuKEoiDihpAg4oSWPC90ZXh0Pjwvc3ZnPg==';

  return (
    <Container>
      <Breadcrumbs items={breadcrumbsItems} />
      
      <Button variant="outline-primary" onClick={handleBackClick} className="mb-3">
        ← Назад к списку устройств
      </Button>

      <Row>
        <Col md={6}>
          <Card>
            <Card.Img 
              variant="top" 
              src={imageUrl}
              style={{ height: '400px', objectFit: 'contain', padding: '2rem' }}
              onError={(e) => {
                e.target.src = 'data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iNDAwIiBoZWlnaHQ9IjQwMCIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj48cmVjdCB3aWR0aD0iMTAwJSIgaGVpZ2h0PSIxMDAlIiBmaWxsPSIjZGRkIi8+PHRleHQgeD0iNTAlIiB5PSI1MCUiIGZvbnQtc2l6ZT0iMjQiIGZpbGw9IiM5OTkiIGRvbWluYW50LWJhc2VsaW5lPSJtaWRkbGUiIHRleHQtYW5jaG9yPSJtaWRkbGUiPuKEoiDihpAg4oSWPC90ZXh0Pjwvc3ZnPg==';
              }}
            />
          </Card>
        </Col>
        
        <Col md={6}>
          <Card>
            <Card.Body>
              <Card.Title as="h2">{device.name}</Card.Title>
              <Card.Subtitle className="mb-3 text-muted">
                {device.model}
              </Card.Subtitle>
              
              <div className="mb-4">
                <h5>Характеристики</h5>
                <Row>
                  <Col sm={6}>
                    <strong>Протокол:</strong>
                    <br />
                    {device.protocol}
                  </Col>
                  <Col sm={6}>
                    <strong>Скорость данных:</strong>
                    <br />
                    {device.avg_data_rate} Кбит/с
                  </Col>
                </Row>
                <Row className="mt-2">
                  <Col sm={6}>
                    <strong>Трафик в час:</strong>
                    <br />
                    {device.data_per_hour} Кб/ч
                  </Col>
                  <Col sm={6}>
                    <strong>Статус:</strong>
                    <br />
                    {device.is_active ? 'Активно' : 'Неактивно'}
                  </Col>
                </Row>
              </div>

              <div className="mb-4">
                <h5>Описание</h5>
                <p>{device.description_all}</p>
              </div>

              <Button variant="primary" size="lg" className="w-100">
                Добавить в корзину
              </Button>
            </Card.Body>
          </Card>
        </Col>
      </Row>
    </Container>
  );
}

export default DeviceDetailPage;