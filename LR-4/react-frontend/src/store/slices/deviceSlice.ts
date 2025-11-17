import { createSlice, type PayloadAction } from '@reduxjs/toolkit'

interface DeviceFilters {
  search: string
  protocol: string
  minDataRate: string
  maxDataRate: string
}

interface DeviceState {
  filters: DeviceFilters
}

const initialState: DeviceState = {
  filters: {
    search: '',
    protocol: '',
    minDataRate: '',
    maxDataRate: '',
  },
}

const deviceSlice = createSlice({
  name: 'devices',
  initialState,
  reducers: {
    setSearchFilter: (state, action: PayloadAction<string>) => {
      state.filters.search = action.payload
    },
    setProtocolFilter: (state, action: PayloadAction<string>) => {
      state.filters.protocol = action.payload
    },
    setMinDataRateFilter: (state, action: PayloadAction<string>) => {
      state.filters.minDataRate = action.payload
    },
    setMaxDataRateFilter: (state, action: PayloadAction<string>) => {
      state.filters.maxDataRate = action.payload
    },
    clearFilters: (state) => {
      state.filters = initialState.filters
    },
  },
})

export const {
  setSearchFilter,
  setProtocolFilter,
  setMinDataRateFilter,
  setMaxDataRateFilter,
  clearFilters,
} = deviceSlice.actions

export default deviceSlice.reducer