import React, { useState, useEffect } from 'react';
import { Container, Row, Col, Form } from 'react-bootstrap';
import Breadcrumbs from '../components/Layout/Breadcrumbs';
import DeviceCard from '../components/Devices/DeviceCard';
import { apiService } from '../services/api';

function DevicesPage() {
  const [devices, setDevices] = useState([]);
  const [loading, setLoading] = useState(true);
  const [filters, setFilters] = useState({
    search: '',
    protocol: ''
  });

  const breadcrumbsItems = [
    { label: 'Умные устройства', active: true }
  ];

  useEffect(() => {
    loadDevices();
  }, [filters]);

  const loadDevices = async () => {
    setLoading(true);
    try {
      const devicesData = await apiService.getDevices(filters);
      setDevices(devicesData);
    } catch (error) {
      console.error('Error loading devices:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleSearchChange = (e) => {
    setFilters({ ...filters, search: e.target.value });
  };

  const handleProtocolChange = (e) => {
    setFilters({ ...filters, protocol: e.target.value });
  };

  return (
    <Container>
      <Breadcrumbs items={breadcrumbsItems} />
      
      <h1>Каталог умных устройств</h1>
      <p className="text-muted mb-4">
        Выберите подходящие устройства для вашего умного дома
      </p>

      {/* Фильтры */}
      <Row className="mb-4">
        <Col md={6}>
          <Form.Group>
            <Form.Label>Поиск устройств</Form.Label>
            <Form.Control
              type="text"
              placeholder="Введите название или описание..."
              value={filters.search}
              onChange={handleSearchChange}
            />
          </Form.Group>
        </Col>
        
        <Col md={6}>
          <Form.Group>
            <Form.Label>Протокол</Form.Label>
            <Form.Select value={filters.protocol} onChange={handleProtocolChange}>
              <option value="">Все протоколы</option>
              <option value="Wi-Fi">Wi-Fi</option>
              <option value="Zigbee">Zigbee</option>
              <option value="Bluetooth">Bluetooth</option>
            </Form.Select>
          </Form.Group>
        </Col>
      </Row>

      {/* Список устройств */}
      {loading ? (
        <div className="text-center py-5">
          <div className="spinner-border" role="status">
            <span className="visually-hidden">Загрузка...</span>
          </div>
        </div>
      ) : devices.length === 0 ? (
        <div className="alert alert-info text-center">
          Устройства не найдены. Попробуйте изменить параметры поиска.
        </div>
      ) : (
        <Row>
          {devices.map(device => (
            <Col key={device.id} xs={12} sm={6} lg={4} className="mb-4">
              <DeviceCard device={device} />
            </Col>
          ))}
        </Row>
      )}
    </Container>
  );
}

export default DevicesPage;