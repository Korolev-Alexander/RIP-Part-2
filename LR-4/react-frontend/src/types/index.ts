export interface User {
  id: number;
  username: string;
  email: string;
  createdAt: string;
}

export interface Device {
  id: number;
  name: string;
  description: string;
  price: number;
  quantity: number;
  imageUrl: string;
  createdAt: string;
}

export interface OrderItem {
  id: number;
  deviceId: number;
  quantity: number;
  price: number;
}

export interface Order {
  id: number;
  items: OrderItem[];
  status: 'draft' | 'submitted' | 'confirmed' | 'shipped' | 'delivered';
  totalAmount: number;
  createdAt: string;
}

export interface AuthState {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
}

export interface OrderState {
  currentOrder: Order | null;
  orderHistory: Order[];
}