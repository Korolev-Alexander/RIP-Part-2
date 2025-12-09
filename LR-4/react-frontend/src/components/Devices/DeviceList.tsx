import React from 'react';
import { Row, Col, Button } from 'react-bootstrap';
import { useAppDispatch, useAppSelector } from '../../store/hooks';
import { addService } from '../../store/slices/orderSlice';
import type { SmartDevice } from '../../api/Api';
import DeviceCard from './DeviceCard';

interface DeviceListProps {
  devices: SmartDevice[];
}

const DeviceList: React.FC<DeviceListProps> = ({ devices }) => {
  const dispatch = useAppDispatch();
  const isAuthenticated = useAppSelector(state => state.user.isAuthenticated);
  
  const handleAddService = () => {
    // Пример добавления услуги с фиксированными данными
    // В реальном приложении данные могут приходить из формы или API
    dispatch(addService({
      id: Date.now(), // Уникальный ID для примера
      name: "Услуга по установке",
      price: 1500
    }));
  };

  if (devices.length === 0) {
    return (
      <div className="text-center py-5">
        <h4>Устройства не найдены</h4>
        <p className="text-muted">Попробуйте изменить параметры фильтрации</p>
      </div>
    );
  }

  return (
    <div>
      {isAuthenticated && (
        <div className="mb-3">
          <Button variant="primary" onClick={handleAddService}>
            Добавить услугу
          </Button>
        </div>
      )}
      <Row className="g-4">
        {devices.map((device) => (
          <Col key={device.id} xs={12} sm={6} lg={4}>
            <DeviceCard device={device} />
          </Col>
        ))}
      </Row>
    </div>
  );
};

export default DeviceList;