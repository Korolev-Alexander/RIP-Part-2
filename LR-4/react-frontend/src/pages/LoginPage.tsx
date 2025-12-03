import React, { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { Container, Form, Button, Card, Alert, Spinner } from 'react-bootstrap';
import { useAppDispatch, useAppSelector } from '../store/hooks';
import { loginStart, loginSuccess, loginFailure } from '../store/slices/authSlice';
import { api } from '../services/api';

const LoginPage: React.FC = () => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const dispatch = useAppDispatch();
  const navigate = useNavigate();
  const { loading, error } = useAppSelector((state) => state.auth);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    dispatch(loginStart());
    
    try {
      // Вызываем API для авторизации
      const userData = await api.login(username, password);
      
      // Сохраняем информацию о пользователе в localStorage
      localStorage.setItem('user', JSON.stringify({
        username: userData.username,
        isAuthenticated: true,
        token: userData.token
      }));
      
      // Обновляем состояние авторизации в Redux
      dispatch(loginSuccess({
        username: userData.username,
        isAuthenticated: true
      }));
      
      // Перенаправляем на главную страницу
      navigate('/');
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Ошибка авторизации';
      dispatch(loginFailure(errorMessage));
    }
  };

  return (
    <Container className="d-flex align-items-center justify-content-center" style={{ minHeight: '100vh' }}>
      <div className="w-100" style={{ maxWidth: '400px' }}>
        <Card>
          <Card.Body>
            <h2 className="text-center mb-4">Вход</h2>
            {error && <Alert variant="danger">{error}</Alert>}
            <Form onSubmit={handleSubmit}>
              <Form.Group id="username" className="mb-3">
                <Form.Label>Имя пользователя</Form.Label>
                <Form.Control
                  type="text"
                  value={username}
                  onChange={(e) => setUsername(e.target.value)}
                  required
                  disabled={loading}
                />
              </Form.Group>
              <Form.Group id="password" className="mb-3">
                <Form.Label>Пароль</Form.Label>
                <Form.Control
                  type="password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  required
                  disabled={loading}
                />
              </Form.Group>
              <Button className="w-100" type="submit" disabled={loading}>
                {loading ? (
                  <>
                    <Spinner
                      as="span"
                      animation="border"
                      size="sm"
                      role="status"
                      aria-hidden="true"
                    />
                    <span className="ms-2">Вход...</span>
                  </>
                ) : (
                  'Войти'
                )}
              </Button>
            </Form>
          </Card.Body>
        </Card>
        <div className="w-100 text-center mt-2">
          Нет аккаунта? <Link to="/register">Зарегистрироваться</Link>
        </div>
      </div>
    </Container>
  );
};

export default LoginPage;