import React, {useEffect} from 'react';
import {useDispatch, useSelector} from 'react-redux';
import './App.css';
import Header from "./components/Header.js";
import Logs from "./components/Logs";

import {fetchConfigAction} from "./actions/configSlice";
import WebSocketProvider from "./components/WebSocketProvider";

const App = () => {
  const dispatch = useDispatch();
  const {loaded} = useSelector((state) => state.config);

  useEffect(() => {
    if (!loaded) {
      dispatch(fetchConfigAction());
    }
  }, [loaded, dispatch]);

  return (
    <WebSocketProvider>
      <Header/>
      <Logs/>
    </WebSocketProvider>
  );
}

export default App;
