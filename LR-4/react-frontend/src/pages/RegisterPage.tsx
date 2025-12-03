import React, { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { Container, Form, Button, Card, Alert, Spinner } from 'react-bootstrap';
import { useAppDispatch, useAppSelector } from '../store/hooks';
import { registerStart, registerSuccess, registerFailure } from '../store/slices/authSlice';
import { api } from '../services/api';

const RegisterPage: React.FC = () => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const dispatch = useAppDispatch();
  const navigate = useNavigate();
  const { loading, error } = useAppSelector((state) => state.auth);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (password !== confirmPassword) {
      dispatch(registerFailure('Пароли не совпадают'));
      return;
    }
    
    dispatch(registerStart());
    
    try {
      // Вызываем API для регистрации
      const userData = await api.register(username, password);
      
      // Сохраняем информацию о пользователе в localStorage
      localStorage.setItem('user', JSON.stringify({
        username: userData.username,
        isAuthenticated: true,
        token: userData.token
      }));
      
      // Обновляем состояние авторизации в Redux
      dispatch(registerSuccess({
        username: userData.username,
        isAuthenticated: true
      }));
      
      // Перенаправляем на главную страницу
      navigate('/');
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Ошибка регистрации';
      dispatch(registerFailure(errorMessage));
    }
  };

  return (
    <Container className="d-flex align-items-center justify-content-center" style={{ minHeight: '100vh' }}>
      <div className="w-100" style={{ maxWidth: '400px' }}>
        <Card>
          <Card.Body>
            <h2 className="text-center mb-4">Регистрация</h2>
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
              <Form.Group id="confirmPassword" className="mb-3">
                <Form.Label>Подтверждение пароля</Form.Label>
                <Form.Control
                  type="password"
                  value={confirmPassword}
                  onChange={(e) => setConfirmPassword(e.target.value)}
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
                    <span className="ms-2">Регистрация...</span>
                  </>
                ) : (
                  'Зарегистрироваться'
                )}
              </Button>
            </Form>
          </Card.Body>
        </Card>
        <div className="w-100 text-center mt-2">
          Уже есть аккаунт? <Link to="/login">Войти</Link>
        </div>
      </div>
    </Container>
  );
};

export default RegisterPage;