import React from 'react';
import { Row, Col } from 'react-bootstrap';
import type { SmartDevice } from '../../types';
import DeviceCard from './DeviceCard';

interface DeviceListProps {
  devices: SmartDevice[];
}

const DeviceList: React.FC<DeviceListProps> = ({ devices }) => {
  if (devices.length === 0) {
    return (
      <div className="text-center py-5">
        <h4>Устройства не найдены</h4>
        <p className="text-muted">Попробуйте изменить параметры фильтрации</p>
      </div>
    );
  }

  return (
    <Row className="g-4">
      {devices.map((device) => (
        <Col key={device.id} xs={12} sm={6} lg={4}>
          <DeviceCard device={device} />
        </Col>
      ))}
    </Row>
  );
};

export default DeviceList;