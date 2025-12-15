import React, { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Container, Table, Button, Alert, Spinner, Card } from 'react-bootstrap';
import { useAppSelector, useAppDispatch } from '../store/hooks';
import { fetchUserOrders } from '../store/slices/orderSlice';
import type { RootState } from '../store/index';
import type { SmartOrder } from '../api/Api';

const OrdersPage: React.FC = () => {
  const navigate = useNavigate();
  const dispatch = useAppDispatch();
  const user = useAppSelector((state: RootState) => state.user);
  const { userOrders, loading, error } = useAppSelector((state: RootState) => state.order);

  useEffect(() => {
    // Проверяем авторизацию пользователя
    if (!user || !user.isAuthenticated) {
      navigate('/login');
      return;
    }

    // Загружаем заявки пользователя только если пользователь авторизован
    if (user.isAuthenticated) {
      dispatch(fetchUserOrders());
    }
  }, [dispatch, user, navigate]);

  const handleViewOrder = (orderId: number) => {
    navigate(`/orders/${orderId}`);
  };

  // Функция для отображения статуса заявки на русском языке
  const getOrderStatusText = (status: string) => {
    switch (status) {
      case 'draft':
        return 'Черновик';
      case 'formed':
        return 'Сформирована';
      case 'completed':
        return 'Завершена';
      case 'rejected':
        return 'Отклонена';
      case 'deleted':
        return 'Удалена';
      default:
        return status;
    }
  };

  // Функция для получения класса статуса для стилизации
  const getStatusClass = (status: string) => {
    switch (status) {
      case 'draft':
        return 'warning';
      case 'formed':
        return 'info';
      case 'completed':
        return 'success';
      case 'rejected':
        return 'danger';
      case 'deleted':
        return 'secondary';
      default:
        return 'secondary';
    }
  };

  return (
    <Container className="mt-4">
      <h2 className="mb-4">Мои заявки</h2>
      
      {error && (
        <Alert variant="danger" className="mb-4">
          Ошибка: {error}
        </Alert>
      )}
      
      {loading ? (
        <div className="text-center">
          <Spinner animation="border" role="status">
            <span className="visually-hidden">Загрузка...</span>
          </Spinner>
        </div>
      ) : (
        <>
          {!userOrders || userOrders.length === 0 ? (
            <Card>
              <Card.Body>
                <Card.Text className="text-center">
                  У вас пока нет заявок. Перейдите в каталог устройств, чтобы создать новую заявку.
                </Card.Text>
                <div className="text-center">
                  <Button
                    variant="primary"
                    onClick={() => navigate('/devices')}
                    className="mt-2"
                  >
                    Перейти к устройствам
                  </Button>
                </div>
              </Card.Body>
            </Card>
          ) : (
            <Table striped bordered hover responsive>
              <thead>
                <tr>
                  <th>Номер заявки</th>
                  <th>Дата создания</th>
                  <th>Адрес</th>
                  <th>Трафик (МБ/мес)</th>
                  <th>Статус</th>
                  <th>Действия</th>
                </tr>
              </thead>
              <tbody>
                {userOrders.map((order: SmartOrder) => (
                  <tr key={order.id}>
                    <td>{order.id}</td>
                    <td>{new Date(order.created_at || '').toLocaleDateString('ru-RU')}</td>
                    <td>{order.address}</td>
                    <td>{order.total_traffic?.toFixed(2)}</td>
                    <td>
                      <span className={`badge bg-${getStatusClass(order.status || '')}`}>
                        {getOrderStatusText(order.status || '')}
                      </span>
                    </td>
                    <td>
                      <Button
                        variant="outline-primary"
                        size="sm"
                        onClick={() => handleViewOrder(order.id!)}
                      >
                        Просмотр
                      </Button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </Table>
          )}
        </>
      )}
    </Container>
  );
};

export default OrdersPage;
