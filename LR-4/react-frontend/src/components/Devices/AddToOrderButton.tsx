import React from 'react';
import { Button, Spinner } from 'react-bootstrap';
import { useAppSelector, useAppDispatch } from '../../store/hooks';
import { addDeviceToDraft, createDraftOrder } from '../../store/slices/orderSlice';
import type { SmartDevice } from '../../types';

interface AddToOrderButtonProps {
  device: SmartDevice;
}

const AddToOrderButton: React.FC<AddToOrderButtonProps> = ({ device }) => {
  const dispatch = useAppDispatch();
  const { user } = useAppSelector((state) => state.auth);
  const { draftOrder, loading } = useAppSelector((state) => state.order);

  const handleAddToOrder = () => {
    // Проверяем, авторизован ли пользователь
    if (!user || !user.isAuthenticated) {
      alert('Пожалуйста, авторизуйтесь для добавления устройств в заявку');
      return;
    }

    // Если нет черновика, создаем его
    if (!draftOrder) {
      dispatch(createDraftOrder({ clientId: 1 })); // В реальной реализации здесь будет реальный ID клиента
    }

    // Добавляем устройство в черновик
    dispatch(addDeviceToDraft({ device, quantity: 1 }));
  };

  // Проверяем, есть ли устройство в черновике
  const isInDraft = draftOrder?.items.some(item => item.device_id === device.id);

  return (
    <Button
      variant={isInDraft ? "success" : "outline-primary"}
      size="sm"
      onClick={handleAddToOrder}
      disabled={loading}
    >
      {loading ? (
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
      ) : isInDraft ? (
        "✓ В заявке"
      ) : (
        "Добавить в заявку"
      )}
    </Button>
  );
};

export default AddToOrderButton;