import {configureStore, getDefaultMiddleware} from '@reduxjs/toolkit';
import {createLogger} from 'redux-logger';

import rootReducer from './reducers';

const middlewares = [];
if (process.env.NODE_ENV === 'development') {
  const loggerMiddleware = createLogger();
  middlewares.push(loggerMiddleware);
}

const store = configureStore({
  reducer: rootReducer,
  middleware: [...middlewares, ...getDefaultMiddleware()]
});

export default store;
