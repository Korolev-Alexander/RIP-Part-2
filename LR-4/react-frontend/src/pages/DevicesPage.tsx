import React, { useState, useEffect } from 'react';
import { Container, Row, Col, Form, Button, Spinner, Alert, Card } from 'react-bootstrap';
import type { SmartDevice, DeviceFilter } from '../types';
import { api } from '../services/api';
import DeviceList from '../components/Devices/DeviceList';

const DevicesPage: React.FC = () => {
  const [devices, setDevices] = useState<SmartDevice[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [filters, setFilters] = useState<DeviceFilter>({});

  useEffect(() => {
    loadDevices();
  }, []);

  const loadDevices = async (filterParams?: DeviceFilter) => {
    try {
      setLoading(true);
      setError(null);
      const devicesData = await api.getDevices(filterParams);
      setDevices(devicesData);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load devices');
    } finally {
      setLoading(false);
    }
  };

  const handleFilterChange = (key: keyof DeviceFilter, value: string) => {
    const newFilters = { ...filters, [key]: value || undefined };
    setFilters(newFilters);
  };

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    loadDevices(filters);
  };

  const handleReset = () => {
    setFilters({});
    loadDevices({});
  };

  const protocols = ['Wi-Fi', 'Bluetooth', 'Zigbee'];

  return (
    <Container className="mt-4">
      <Row>
        <Col>
          <h1 className="mb-4">Умные устройства</h1>
          
          {/* Фильтры */}
          <Card className="mb-4">
            <Card.Body>
              <Form onSubmit={handleSearch}>
                <Row className="g-3">
                  <Col md={6}>
                    <Form.Group>
                      <Form.Label>Поиск по названию</Form.Label>
                      <Form.Control
                        type="text"
                        placeholder="Введите название устройства..."
                        value={filters.search || ''}
                        onChange={(e) => handleFilterChange('search', e.target.value)}
                      />
                    </Form.Group>
                  </Col>
                  
                  <Col md={4}>
                    <Form.Group>
                      <Form.Label>Протокол</Form.Label>
                      <Form.Select
                        value={filters.protocol || ''}
                        onChange={(e) => handleFilterChange('protocol', e.target.value)}
                      >
                        <option value="">Все протоколы</option>
                        {protocols.map(protocol => (
                          <option key={protocol} value={protocol}>
                            {protocol}
                          </option>
                        ))}
                      </Form.Select>
                    </Form.Group>
                  </Col>
                  
                  <Col md={2} className="d-flex align-items-end">
                    <div className="d-grid gap-2 w-100">
                      <Button type="submit" variant="primary">
                        Применить
                      </Button>
                      <Button 
                        type="button" 
                        variant="outline-secondary"
                        onClick={handleReset}
                      >
                        Сбросить
                      </Button>
                    </div>
                  </Col>
                </Row>
              </Form>
            </Card.Body>
          </Card>

          {/* Результаты */}
          {error && (
            <Alert variant="danger" className="mb-4">
              {error}
            </Alert>
          )}

          {loading ? (
            <div className="text-center">
              <Spinner animation="border" role="status">
                <span className="visually-hidden">Загрузка...</span>
              </Spinner>
            </div>
          ) : (
            <DeviceList devices={devices} />
          )}
        </Col>
      </Row>
    </Container>
  );
};

export default DevicesPage;