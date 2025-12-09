import React from 'react';
import { Button, Spinner } from 'react-bootstrap';
import { useAppSelector, useAppDispatch } from '../../store/hooks';
import { addItem } from '../../store/slices/orderSlice';
import type { SmartDevice } from '../../api/Api';

interface AddToOrderButtonProps {
  device: SmartDevice;
}

const AddToOrderButton: React.FC<AddToOrderButtonProps> = ({ device }) => {
  const dispatch = useAppDispatch();
  const orderState = useAppSelector((state) => state.order);

  const handleAddToOrder = () => {
    // Добавляем устройство в заявку
    dispatch(addItem({
      id: Date.now(), // Генерируем временный ID для элемента
      deviceId: device.id!,
      quantity: 1,
      price: device.avg_data_rate || 0
    }));
  };

  // Проверяем, есть ли устройство в заявке
  const isInOrder = orderState.items.some(item => item.deviceId === device.id);

  return (
    <Button
      variant={isInOrder ? "success" : "outline-primary"}
      size="sm"
      onClick={handleAddToOrder}
      disabled={orderState.loading}
    >
      {orderState.loading ? (
        <>
          <Spinner
            as="span"
            animation="border"
            size="sm"
            role="status"
            aria-hidden="true"
          />
          <span className="ms-1">Добавление...</span>
        </>
      ) : isInOrder ? (
        "✓ В заявке"
      ) : (
        "Добавить в заявку"
      )}
    </Button>
  );
};

export default AddToOrderButton;