import React, { useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { Container, Card, ListGroup, Badge, Button, Alert, Spinner } from 'react-bootstrap';
import { useAppSelector, useAppDispatch } from '../store/hooks';
import { fetchUserOrders } from '../store/slices/orderSlice';
import type { RootState } from '../store/index';
import type { SmartOrder, OrderItem as ApiOrderItem } from '../api/Api';

// Определяем типы для состояния пользователя и заявок
interface UserState {
  id: number | null;
  username: string | null;
  email: string | null;
  isAuthenticated: boolean;
  token: string | null;
}

interface OrderState {
  id: number | null;
  items: any[];
  status: 'draft' | 'submitted' | 'confirmed' | 'shipped' | 'delivered';
  totalAmount: number;
  createdAt: string | null;
  loading: boolean;
  error: string | null;
  userOrders: SmartOrder[];
}

const OrderPage: React.FC = () => {
  const { id } = useParams<{ id?: string }>();
  const navigate = useNavigate();
  const dispatch = useAppDispatch();
  const user = useAppSelector((state: RootState) => state.user) as UserState;
  const orderState = useAppSelector((state: RootState) => state.order) as OrderState;

  // Находим заявку по ID
  const order = orderState.userOrders.find((o: SmartOrder) => o.id === Number(id));

  useEffect(() => {
    // Проверяем авторизацию пользователя
    if (!user || !user.isAuthenticated) {
      navigate('/login');
      return;
    }

    // Если заявка не найдена в хранилище или список заявок пуст, загружаем заявки
    if (!order && orderState.userOrders.length === 0) {
      dispatch(fetchUserOrders());
    }
  }, [dispatch, user, order, orderState.userOrders, navigate]);

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

  if (orderState.loading) {
    return (
      <Container className="mt-4 text-center">
        <Spinner animation="border" role="status">
          <span className="visually-hidden">Загрузка...</span>
        </Spinner>
      </Container>
    );
  }

  if (orderState.error) {
    return (
      <Container className="mt-4">
        <Alert variant="danger">Ошибка: {orderState.error}</Alert>
        <Button variant="primary" onClick={() => navigate('/orders')}>
          Вернуться к списку заявок
        </Button>
      </Container>
    );
  }

  if (!order) {
    return (
      <Container className="mt-4">
        <Alert variant="warning">Заявка не найдена</Alert>
        <Button variant="primary" onClick={() => navigate('/orders')}>
          Вернуться к списку заявок
        </Button>
      </Container>
    );
  }

  return (
    <Container className="mt-4">
      <div className="d-flex justify-content-between align-items-center mb-4">
        <h2>Детали заявки #{order.id}</h2>
        <Button variant="secondary" onClick={() => navigate('/orders')}>
          Назад к списку
        </Button>
      </div>

      <Card className="mb-4">
        <Card.Body>
          <div className="d-flex justify-content-between">
            <div>
              <h5>Информация о заявке</h5>
              <p className="mb-1">
                <strong>Статус:</strong>{' '}
                <Badge bg={getStatusClass(order.status || '')}>
                  {getOrderStatusText(order.status || '')}
                </Badge>
              </p>
              <p className="mb-1">
                <strong>Адрес:</strong> {order.address}
              </p>
              <p className="mb-1">
                <strong>Дата создания:</strong>{' '}
                {new Date(order.created_at || '').toLocaleDateString('ru-RU')}
              </p>
              {order.formed_at && (
                <p className="mb-1">
                  <strong>Дата формирования:</strong>{' '}
                  {new Date(order.formed_at).toLocaleDateString('ru-RU')}
                </p>
              )}
              <p className="mb-1">
                <strong>Трафик:</strong> {order.total_traffic?.toFixed(2)} МБ/мес
              </p>
            </div>
          </div>
        </Card.Body>
      </Card>

      <Card>
        <Card.Header>
          <h5 className="mb-0">Устройства в заявке</h5>
        </Card.Header>
        <ListGroup variant="flush">
          {order.items && order.items.length > 0 ? (
            order.items.map((item: ApiOrderItem) => (
              <ListGroup.Item key={item.device_id}>
                <div className="d-flex justify-content-between align-items-center">
                  <div>
                    <h6 className="mb-1">{item.device_name}</h6>
                    <small className="text-muted">
                      Трафик: {item.data_per_hour?.toFixed(2)} МБ/час
                    </small>
                  </div>
                  <div className="text-end">
                    <div>Количество: {item.quantity}</div>
                  </div>
                </div>
              </ListGroup.Item>
            ))
          ) : (
            <ListGroup.Item>В заявке нет устройств</ListGroup.Item>
          )}
        </ListGroup>
      </Card>
    </Container>
  );
};

export default OrderPage;