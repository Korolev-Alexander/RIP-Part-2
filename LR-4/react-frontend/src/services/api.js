// Mock данные для случаев, когда бэкенд недоступен
const mockDevices = [
  {
    id: 1,
    name: "Хаб",
    model: "Яндекс Хаб",
    avg_data_rate: 5120,
    data_per_hour: 56.25,
    namespace_url: "http://localhost:9000/image/hub.png",
    description: "Умный пульт Яндекс Хаб для устройств",
    description_all: "Умный пульт Яндекс Хаб для управления всеми устройствами умного дома...",
    protocol: "Wi-Fi",
    is_active: true,
    created_at: "2025-10-21T13:08:04Z"
  },
  {
    id: 2,
    name: "Лампочка",
    model: "Яндекс, E27",
    avg_data_rate: 8,
    data_per_hour: 0.5,
    namespace_url: "http://localhost:9000/image/lamp.png",
    description: "Умная лампочка Яндекс, E27",
    description_all: "Умная Яндекс лампочка позволяет дистанционно управлять освещением...",
    protocol: "Wi-Fi",
    is_active: true,
    created_at: "2025-10-21T13:08:04Z"
  },
  {
    id: 3,
    name: "Розетка",
    model: "YNDX-00340",
    avg_data_rate: 2,
    data_per_hour: 0.1,
    namespace_url: "",
    description: "Умная розетка Яндекс YNDX-00340",
    description_all: "Умная розетка для дистанционного управления электроприборами...",
    protocol: "Wi-Fi",
    is_active: true,
    created_at: "2025-10-21T13:08:04Z"
  }
];

class ApiService {
  async getDevices(filters = {}) {
    try {
      const queryParams = new URLSearchParams();
      if (filters.search) queryParams.append('search', filters.search);
      if (filters.protocol) queryParams.append('protocol', filters.protocol);

      const response = await fetch(`/api/smart-devices?${queryParams}`);
      
      if (!response.ok) throw new Error('API error');
      
      return await response.json();
    } catch (error) {
      console.error('API error, using mock data:', error);
      // Фильтрация mock данных
      return this.filterMockDevices(mockDevices, filters);
    }
  }

  async getDevice(id) {
    try {
      const response = await fetch(`/api/smart-devices/${id}`);
      
      if (!response.ok) throw new Error('API error');
      
      return await response.json();
    } catch (error) {
      console.error('API error, using mock data:', error);
      return mockDevices.find(device => device.id === parseInt(id)) || null;
    }
  }

  filterMockDevices(devices, filters) {
    let filtered = devices.filter(device => device.is_active);

    if (filters.search) {
      const searchLower = filters.search.toLowerCase();
      filtered = filtered.filter(device =>
        device.name.toLowerCase().includes(searchLower) ||
        device.description.toLowerCase().includes(searchLower)
      );
    }

    if (filters.protocol) {
      filtered = filtered.filter(device => device.protocol === filters.protocol);
    }

    return filtered;
  }
}

export const apiService = new ApiService();