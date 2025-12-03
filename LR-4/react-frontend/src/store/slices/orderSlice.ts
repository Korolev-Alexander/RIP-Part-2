import { createSlice, type PayloadAction } from '@reduxjs/toolkit';
import type { SmartOrder, SmartOrderItem, SmartDevice } from '../../types';

interface OrderState {
  draftOrder: SmartOrder | null;
  orders: SmartOrder[];
  loading: boolean;
  error: string | null;
}

const initialState: OrderState = {
  draftOrder: null,
  orders: [],
  loading: false,
  error: null,
};

const orderSlice = createSlice({
  name: 'order',
  initialState,
  reducers: {
    // Создание черновика заявки
    createDraftOrder: (state, action: PayloadAction<{ clientId: number }>) => {
      const now = new Date().toISOString();
      state.draftOrder = {
        id: 0, // Временный ID до сохранения в БД
        status: 'draft',
        address: '',
        total_traffic: 0,
        client_id: action.payload.clientId,
        client_name: '',
        created_at: now,
        items: [],
      };
    },
    
    // Добавление устройства в черновик заявки
    addDeviceToDraft: (state, action: PayloadAction<{ device: SmartDevice; quantity: number }>) => {
      if (!state.draftOrder) return;
      
      const { device, quantity } = action.payload;
      const existingItem = state.draftOrder.items.find(item => item.device_id === device.id);
      
      if (existingItem) {
        existingItem.quantity += quantity;
      } else {
        state.draftOrder.items.push({
          device_id: device.id,
          device_name: device.name,
          quantity,
          data_per_hour: device.data_per_hour,
          namespace_url: device.namespace_url,
        });
      }
      
      // Пересчитываем общий трафик
      state.draftOrder.total_traffic = state.draftOrder.items.reduce(
        (sum, item) => sum + (item.data_per_hour * item.quantity), 
        0
      );
    },
    
    // Изменение количества устройства в черновике
    updateDeviceQuantity: (state, action: PayloadAction<{ deviceId: number; quantity: number }>) => {
      if (!state.draftOrder) return;
      
      const { deviceId, quantity } = action.payload;
      const item = state.draftOrder.items.find(item => item.device_id === deviceId);
      
      if (item) {
        if (quantity <= 0) {
          // Удаляем элемент, если количество <= 0
          state.draftOrder.items = state.draftOrder.items.filter(item => item.device_id !== deviceId);
        } else {
          item.quantity = quantity;
        }
        
        // Пересчитываем общий трафик
        state.draftOrder.total_traffic = state.draftOrder.items.reduce(
          (sum, item) => sum + (item.data_per_hour * item.quantity), 
          0
        );
      }
    },
    
    // Удаление устройства из черновика
    removeDeviceFromDraft: (state, action: PayloadAction<number>) => {
      if (!state.draftOrder) return;
      
      const deviceId = action.payload;
      state.draftOrder.items = state.draftOrder.items.filter(item => item.device_id !== deviceId);
      
      // Пересчитываем общий трафик
      state.draftOrder.total_traffic = state.draftOrder.items.reduce(
        (sum, item) => sum + (item.data_per_hour * item.quantity), 
        0
      );
    },
    
    // Очистка черновика
    clearDraftOrder: (state) => {
      state.draftOrder = null;
    },
    
    // Загрузка заявок пользователя
    loadOrdersStart: (state) => {
      state.loading = true;
      state.error = null;
    },
    
    loadOrdersSuccess: (state, action: PayloadAction<SmartOrder[]>) => {
      state.loading = false;
      state.orders = action.payload;
    },
    
    loadOrdersFailure: (state, action: PayloadAction<string>) => {
      state.loading = false;
      state.error = action.payload;
    },
  },
});

export const {
  createDraftOrder,
  addDeviceToDraft,
  updateDeviceQuantity,
  removeDeviceFromDraft,
  clearDraftOrder,
  loadOrdersStart,
  loadOrdersSuccess,
  loadOrdersFailure,
} = orderSlice.actions;

export default orderSlice.reducer;