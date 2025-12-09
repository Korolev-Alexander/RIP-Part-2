import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Form, Button, Container, Row, Col, Alert } from 'react-bootstrap';
import api from '../api';

const RegisterPage: React.FC = () => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [error, setError] = useState('');
  const [success, setSuccess] = useState(false);
  const [loading, setLoading] = useState(false);
  
  const navigate = useNavigate();

  const validateForm = () => {
    if (!username.trim()) {
      setError('Имя пользователя обязательно');
      return false;
    }
    if (username.length < 3) {
      setError('Имя пользователя должно содержать не менее 3 символов');
      return false;
    }
    if (!password) {
      setError('Пароль обязателен');
      return false;
    }
    if (password.length < 6) {
      setError('Пароль должен содержать не менее 6 символов');
      return false;
    }
    if (password !== confirmPassword) {
      setError('Пароли не совпадают');
      return false;
    }
    return true;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setSuccess(false);
    
    if (!validateForm()) {
      return;
    }
    
    setLoading(true);
    
    try {
      const response = await api.clients.registerCreate({ username, password });
      
      if (response.data.id) {
        setSuccess(true);
        // Автоматически перенаправляем на страницу входа через 2 секунды
        setTimeout(() => {
          navigate('/login');
        }, 2000);
      } else {
        setError('Ошибка регистрации');
      }
    } catch (err: any) {
      console.error('Ошибка регистрации:', err);
      if (err.response?.data?.error) {
        setError(err.response.data.error);
      } else if (err.response?.data?.message) {
        setError(err.response.data.message);
      } else {
        setError('Ошибка подключения к серверу');
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <Container className="mt-5">
      <Row className="justify-content-md-center">
        <Col md={6}>
          <h2 className="text-center mb-4">Регистрация</h2>
          
          {error && <Alert variant="danger">{error}</Alert>}
          {success && <Alert variant="success">Регистрация успешна! Перенаправление на страницу входа...</Alert>}
          
          <Form onSubmit={handleSubmit}>
            <Form.Group className="mb-3" controlId="formUsername">
              <Form.Label>Имя пользователя</Form.Label>
              <Form.Control
                type="text"
                placeholder="Введите имя пользователя"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                disabled={loading}
              />
            </Form.Group>
            
            <Form.Group className="mb-3" controlId="formPassword">
              <Form.Label>Пароль</Form.Label>
              <Form.Control
                type="password"
                placeholder="Введите пароль"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                disabled={loading}
              />
            </Form.Group>
            
            <Form.Group className="mb-3" controlId="formConfirmPassword">
              <Form.Label>Подтверждение пароля</Form.Label>
              <Form.Control
                type="password"
                placeholder="Подтвердите пароль"
                value={confirmPassword}
                onChange={(e) => setConfirmPassword(e.target.value)}
                disabled={loading}
              />
            </Form.Group>
            
            <div className="d-grid">
              <Button 
                variant="primary" 
                type="submit" 
                disabled={loading}
              >
                {loading ? 'Регистрация...' : 'Зарегистрироваться'}
              </Button>
            </div>
          </Form>
          
          <div className="text-center mt-3">
            <p>
              Уже есть аккаунт? <a href="/login">Войдите</a>
            </p>
          </div>
        </Col>
      </Row>
    </Container>
  );
};

export default RegisterPage;