import {combineReducers} from 'redux';
import mainSliceReducer from './actions/mainSlice';
import configSliceReducer from './actions/configSlice';
import wsSliceReducer from './actions/wsSlice';

const appReducer = combineReducers({
  main: mainSliceReducer,
  config: configSliceReducer,
  ws: wsSliceReducer
});

const rootReducer = (state, action) => {
  return appReducer(state, action);
};

export default rootReducer;
