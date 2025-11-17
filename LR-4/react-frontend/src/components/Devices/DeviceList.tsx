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
        <p className="text-muted">Попробуйте изменить параметры поиска</p>
      </div>
    );
  }

  return (
    <Row className="g-3">
      {devices.map((device) => (
        <Col 
          key={device.id} 
          xs={12}     // На мобильных - 1 колонка
          sm={6}      // На планшетах - 2 колонки  
          md={4}      // На десктопах - 3 колонки
          lg={3}      // На больших экранах - 4 колонки
        >
          <DeviceCard device={device} />
        </Col>
      ))}
    </Row>
  );
};

export default DeviceList;