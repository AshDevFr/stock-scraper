import {createAsyncThunk, createSlice} from '@reduxjs/toolkit';

import {fetchConfig} from '../services/apiClient';

export const fetchConfigAction = createAsyncThunk('config', fetchConfig);

const configSlice = createSlice({
  name: 'config',
  initialState: {
    config: {},
    loaded: false
  },
  reducers: {},
  extraReducers: {
    [fetchConfigAction.pending]: (state) => {
      state.status = 'loading';
    },
    [fetchConfigAction.fulfilled]: (state, {payload: {data}}) => {
      state.config = data;
      state.loaded = true
      state.status = 'succeeded';
    },
    [fetchConfigAction.rejected]: (state) => {
      state.status = 'failed';
    }
  }
});

export default configSlice.reducer;
