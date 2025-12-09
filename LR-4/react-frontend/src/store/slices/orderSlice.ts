import { createSlice, createAsyncThunk, type PayloadAction } from '@reduxjs/toolkit';
import api from '../../api';
import type { SmartOrder } from '../../api/Api';

interface OrderItem {
  id: number;
  deviceId: number;
  quantity: number;
  price: number;
}

interface Service {
  id: number;
  name: string;
  price: number;
}

interface OrderState {
  id: number | null;
  items: OrderItem[];
  services: Service[];
  status: 'draft' | 'submitted' | 'confirmed' | 'shipped' | 'delivered';
  totalAmount: number;
  createdAt: string | null;
  loading: boolean;
  error: string | null;
  userOrders: SmartOrder[];
}

const initialState: OrderState = {
  id: null,
  items: [],
  services: [],
  status: 'draft',
  totalAmount: 0,
  createdAt: null,
  loading: false,
  error: null,
  userOrders: [],
};

// Async Thunk функции
export const fetchUserOrders = createAsyncThunk(
  'order/fetchUserOrders',
  async (_, { rejectWithValue }) => {
    try {
      const response = await api.smartOrders.smartOrdersList();
      return response.data;
    } catch (error: any) {
      return rejectWithValue(error.response?.data?.message || 'Ошибка при загрузке заявок');
    }
  }
);

export const createOrder = createAsyncThunk(
  'order/createOrder',
  async (orderData: { address: string }, { rejectWithValue }) => {
    try {
      // Сначала создаем заявку с адресом
      const orderResponse = await api.smartOrders.smartOrdersUpdate(0, {
        address: orderData.address,
      });
      
      // Затем формируем заявку (переводим из draft в formed)
      const formedResponse = await api.smartOrders.formUpdate(orderResponse.data.id!);
      
      return formedResponse.data;
    } catch (error: any) {
      return rejectWithValue(error.response?.data?.message || 'Ошибка при создании заявки');
    }
  }
);

export const updateOrder = createAsyncThunk(
  'order/updateOrder',
  async ({ id, data }: { id: number; data: { address: string } }, { rejectWithValue }) => {
    try {
      const response = await api.smartOrders.smartOrdersUpdate(id, data);
      return response.data;
    } catch (error: any) {
      return rejectWithValue(error.response?.data?.message || 'Ошибка при обновлении заявки');
    }
  }
);

export const deleteOrder = createAsyncThunk(
  'order/deleteOrder',
  async (id: number, { rejectWithValue }) => {
    try {
      await api.smartOrders.smartOrdersDelete(id);
      return id;
    } catch (error: any) {
      return rejectWithValue(error.response?.data?.message || 'Ошибка при удалении заявки');
    }
  }
);

export const orderSlice = createSlice({
  name: 'order',
  initialState,
  reducers: {
    addItem: (state, action: PayloadAction<OrderItem>) => {
      const existingItem = state.items.find(item => item.deviceId === action.payload.deviceId);
      if (existingItem) {
        existingItem.quantity += action.payload.quantity;
      } else {
        state.items.push(action.payload);
      }
      state.totalAmount = state.items.reduce((sum, item) => sum + (item.price * item.quantity), 0) +
                         state.services.reduce((sum, service) => sum + service.price, 0);
    },
    removeItem: (state, action: PayloadAction<number>) => {
      state.items = state.items.filter(item => item.id !== action.payload);
      state.totalAmount = state.items.reduce((sum, item) => sum + (item.price * item.quantity), 0) +
                         state.services.reduce((sum, service) => sum + service.price, 0);
    },
    updateItemQuantity: (state, action: PayloadAction<{ id: number; quantity: number }>) => {
      const item = state.items.find(item => item.id === action.payload.id);
      if (item) {
        item.quantity = action.payload.quantity;
        state.totalAmount = state.items.reduce((sum, item) => sum + (item.price * item.quantity), 0) +
                           state.services.reduce((sum, service) => sum + service.price, 0);
      }
    },
    clearOrder: (state) => {
      state.id = null;
      state.items = [];
      state.services = [];
      state.status = 'draft';
      state.totalAmount = 0;
      state.createdAt = null;
    },
    submitOrder: (state) => {
      state.status = 'submitted';
      state.createdAt = new Date().toISOString();
    },
    addService: (state, action: PayloadAction<Service>) => {
      const existingService = state.services.find(service => service.id === action.payload.id);
      if (!existingService) {
        state.services.push(action.payload);
        state.totalAmount += action.payload.price;
      }
    },
    removeService: (state, action: PayloadAction<number>) => {
      const serviceIndex = state.services.findIndex(service => service.id === action.payload);
      if (serviceIndex !== -1) {
        const service = state.services[serviceIndex];
        state.totalAmount -= service.price;
        state.services.splice(serviceIndex, 1);
      }
    },
  },
  extraReducers: (builder) => {
    builder
      // fetchUserOrders
      .addCase(fetchUserOrders.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchUserOrders.fulfilled, (state, action: PayloadAction<SmartOrder[]>) => {
        state.loading = false;
        state.userOrders = action.payload;
      })
      .addCase(fetchUserOrders.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload as string;
      })
      // createOrder
      .addCase(createOrder.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(createOrder.fulfilled, (state, action: PayloadAction<SmartOrder>) => {
        state.loading = false;
        // Добавляем новую заявку в список заявок пользователя
        state.userOrders.push(action.payload);
      })
      .addCase(createOrder.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload as string;
      })
      // updateOrder
      .addCase(updateOrder.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(updateOrder.fulfilled, (state, action: PayloadAction<SmartOrder>) => {
        state.loading = false;
        // Обновляем заявку в списке заявок пользователя
        const index = state.userOrders.findIndex(order => order.id === action.payload.id);
        if (index !== -1) {
          state.userOrders[index] = action.payload;
        }
      })
      .addCase(updateOrder.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload as string;
      })
      // deleteOrder
      .addCase(deleteOrder.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(deleteOrder.fulfilled, (state, action: PayloadAction<number>) => {
        state.loading = false;
        // Удаляем заявку из списка заявок пользователя
        state.userOrders = state.userOrders.filter(order => order.id !== action.payload);
      })
      .addCase(deleteOrder.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload as string;
      });
  },
});

export const { addItem, removeItem, updateItemQuantity, clearOrder, submitOrder, addService, removeService } = orderSlice.actions;

export default orderSlice.reducer;