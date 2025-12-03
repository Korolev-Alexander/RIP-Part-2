import axios from 'axios';
import type { SmartDevice, SmartOrder, Client, DeviceFilter } from '../types';

const API_BASE_URL = '/api'; // Прокси через Vite

// Создаем экземпляр axios с базовыми настройками
const apiClient = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Интерцептор для добавления токена авторизации
apiClient.interceptors.request.use((config) => {
  const user = localStorage.getItem('user');
  if (user) {
    const userData = JSON.parse(user);
    if (userData.token) {
      config.headers.Authorization = `Bearer ${userData.token}`;
    }
  }
  return config;
});

// Интерцептор для обработки ошибок
apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // Очищаем данные авторизации при ошибке 401
      localStorage.removeItem('user');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

export const api = {
  // ===== AUTH =====
  async login(username: string, password: string) {
    try {
      const response = await apiClient.post('/auth/login', { username, password });
      return response.data;
    } catch (error) {
      throw new Error('Ошибка авторизации');
    }
  },

  async register(username: string, password: string) {
    try {
      const response = await apiClient.post('/auth/register', { username, password });
      return response.data;
    } catch (error) {
      throw new Error('Ошибка регистрации');
    }
  },

  // ===== DEVICES =====
  async getDevices(filters?: DeviceFilter): Promise<SmartDevice[]> {
    try {
      const queryParams = new URLSearchParams();
      if (filters?.search) queryParams.append('search', filters.search);
      if (filters?.protocol) queryParams.append('protocol', filters.protocol);

      const url = `/smart-devices${queryParams.toString() ? `?${queryParams.toString()}` : ''}`;
      const response = await apiClient.get(url);
      return response.data;
    } catch (error) {
      console.error('API error, using mock data:', error);
      // Mock данные для демонстрации
      return [
        {
          id: 1,
          name: 'Умная лампочка',
          model: 'Яндекс, E27',
          avg_data_rate: 8,
          data_per_hour: 0.5,
          namespace_url: '',
          description: 'Умная лампочка Яндекс, E27',
          description_all: 'Умная Яндекс лампочка позволяет дистанционно управлять освещением',
          protocol: 'Wi-Fi',
          is_active: true,
          created_at: new Date().toISOString()
        },
        {
          id: 2,
          name: 'Умная розетка',
          model: 'YNDX-00340',
          avg_data_rate: 2,
          data_per_hour: 0.1,
          namespace_url: '',
          description: 'Умная розетка Яндекс YNDX-00340',
          description_all: 'Умная розетка для дистанционного управления электроприборами',
          protocol: 'Wi-Fi',
          is_active: true,
          created_at: new Date().toISOString()
        }
      ];
    }
  },

  async getDevice(id: number): Promise<SmartDevice> {
    try {
      const response = await apiClient.get(`/smart-devices/${id}`);
      return response.data;
    } catch (error) {
      console.error('API error:', error);
      // Mock данные
      return {
        id,
        name: 'Mock Device',
        model: 'Mock Model',
        avg_data_rate: 10,
        data_per_hour: 1,
        namespace_url: '',
        description: 'Mock description',
        description_all: 'Mock full description',
        protocol: 'Wi-Fi',
        is_active: true,
        created_at: new Date().toISOString()
      };
    }
  },

  // ===== ORDERS =====
  async getOrders(): Promise<SmartOrder[]> {
    try {
      const response = await apiClient.get('/smart-orders');
      return response.data;
    } catch (error) {
      console.error('API error:', error);
      return [];
    }
  },

  async getOrder(id: number): Promise<SmartOrder> {
    try {
      const response = await apiClient.get(`/smart-orders/${id}`);
      return response.data;
    } catch (error) {
      throw new Error('Ошибка загрузки заявки');
    }
  },

  async createOrder(order: Omit<SmartOrder, 'id' | 'created_at'>): Promise<SmartOrder> {
    try {
      const response = await apiClient.post('/smart-orders', order);
      return response.data;
    } catch (error) {
      throw new Error('Ошибка создания заявки');
    }
  },

  async updateOrder(id: number, order: Partial<SmartOrder>): Promise<SmartOrder> {
    try {
      const response = await apiClient.put(`/smart-orders/${id}`, order);
      return response.data;
    } catch (error) {
      throw new Error('Ошибка обновления заявки');
    }
  },

  // ===== CLIENTS =====
  async getClients(): Promise<Client[]> {
    try {
      const response = await apiClient.get('/clients');
      return response.data;
    } catch (error) {
      console.error('API error:', error);
      return [];
    }
  }
};