import React, { useState } from 'react';
import { Container, Form, Button, Card, Alert } from 'react-bootstrap';
import { useAppSelector } from '../store/hooks';

const ProfilePage: React.FC = () => {
  const { user } = useAppSelector((state) => state.auth);
  const [currentPassword, setCurrentPassword] = useState('');
  const [newPassword, setNewPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [message, setMessage] = useState('');
  const [error, setError] = useState('');

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    
    if (newPassword !== confirmPassword) {
      setError('Новые пароли не совпадают');
      return;
    }
    
    if (newPassword.length < 6) {
      setError('Новый пароль должен содержать минимум 6 символов');
      return;
    }
    
    // В реальной реализации здесь будет вызов API для смены пароля
    setMessage('Пароль успешно изменен');
    setCurrentPassword('');
    setNewPassword('');
    setConfirmPassword('');
    
    // Очищаем сообщения через 3 секунды
    setTimeout(() => {
      setMessage('');
    }, 3000);
  };

  return (
    <Container className="mt-4">
      <Card>
        <Card.Header>
          <h3>Личный кабинет</h3>
        </Card.Header>
        <Card.Body>
          <div className="mb-4">
            <h5>Информация о пользователе</h5>
            <p><strong>Имя пользователя:</strong> {user?.username || 'Гость'}</p>
          </div>
          
          <hr />
          
          <h5 className="mt-4">Смена пароля</h5>
          {error && <Alert variant="danger">{error}</Alert>}
          {message && <Alert variant="success">{message}</Alert>}
          
          <Form onSubmit={handleSubmit}>
            <Form.Group className="mb-3">
              <Form.Label>Текущий пароль</Form.Label>
              <Form.Control
                type="password"
                value={currentPassword}
                onChange={(e) => setCurrentPassword(e.target.value)}
                required
              />
            </Form.Group>
            
            <Form.Group className="mb-3">
              <Form.Label>Новый пароль</Form.Label>
              <Form.Control
                type="password"
                value={newPassword}
                onChange={(e) => setNewPassword(e.target.value)}
                required
              />
            </Form.Group>
            
            <Form.Group className="mb-3">
              <Form.Label>Подтверждение нового пароля</Form.Label>
              <Form.Control
                type="password"
                value={confirmPassword}
                onChange={(e) => setConfirmPassword(e.target.value)}
                required
              />
            </Form.Group>
            
            <Button variant="primary" type="submit">
              Изменить пароль
            </Button>
          </Form>
        </Card.Body>
      </Card>
    </Container>
  );
};

export default ProfilePage;