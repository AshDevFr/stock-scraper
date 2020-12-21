import {createSlice} from '@reduxjs/toolkit';

const wsSlice = createSlice({
  name: 'ws',
  initialState: {
    messages: []
  },
  reducers: {
    addMessage(state, action) {
      state.messages.push(action.payload)
    }
  },
  extraReducers: {}
});

export const {addMessage} = wsSlice.actions;

export default wsSlice.reducer;
