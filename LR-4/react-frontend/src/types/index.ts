export interface SmartDevice {
  id: number;
  name: string;
  model: string;
  avg_data_rate: number;
  data_per_hour: number;
  namespace_url: string;
  description: string;
  description_all: string;
  protocol: string;
  is_active: boolean;
  created_at: string;
}

export interface SmartOrder {
  id: number;
  status: string;
  address: string;
  total_traffic: number;
  client_id: number;
  client_name: string;
  formed_at?: string;
  completed_at?: string;
  moderator_id?: number;
  moderator_name?: string;
  created_at: string;
  items: SmartOrderItem[];
}

export interface SmartOrderItem {
  device_id: number;
  device_name: string;
  quantity: number;
  data_per_hour: number;
  namespace_url: string;
}

export interface Client {
  id: number;
  username: string;
  is_moderator: boolean;
  is_active: boolean;
}

export interface DeviceFilter {
  search?: string;
  protocol?: string;
}